package loading

import (
	"container/list"

	"github.com/thewizardplusplus/go-blockchain"
)

// Parameters ...
type Parameters struct {
	Cursor interface{}
	Count  int
}

// Results ...
type Results struct {
	Blocks     blockchain.BlockGroup
	NextCursor interface{}
}

type bucket struct {
	key   Parameters
	value Results
}

type bucketGroup map[Parameters]*list.Element

// LRUCache ...
type LRUCache struct {
	maximalSize int

	buckets bucketGroup
	queue   *list.List
}

// NewLRUCache ...
func NewLRUCache(maximalSize int) LRUCache {
	return LRUCache{
		maximalSize: maximalSize,

		buckets: make(bucketGroup),
		queue:   list.New(),
	}
}

// Get ...
func (cache LRUCache) Get(
	parameters Parameters,
) (results Results, isFound bool) {
	element, isFound := cache.getAndLiftElement(parameters)
	if !isFound {
		return Results{}, false
	}

	return element.Value.(bucket).value, true
}

// Set ...
func (cache LRUCache) Set(parameters Parameters, results Results) {
	newBucket := bucket{parameters, results}
	if element, isFound := cache.getAndLiftElement(parameters); isFound {
		element.Value = newBucket
		return
	}

	// add the new element at the beginning
	element := cache.queue.PushFront(newBucket)
	cache.buckets[parameters] = element
	if cache.queue.Len() <= cache.maximalSize {
		return
	}

	// if the size exceeds the maximum remove the last element
	element = cache.queue.Back()
	cache.queue.Remove(element)
	delete(cache.buckets, element.Value.(bucket).key)
}

func (cache LRUCache) getAndLiftElement(
	parameters Parameters,
) (element *list.Element, isFound bool) {
	element, isFound = cache.buckets[parameters]
	if isFound {
		cache.queue.MoveToFront(element)
	}

	return element, isFound
}
