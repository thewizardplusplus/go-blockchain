package loading

import (
	"github.com/thewizardplusplus/go-blockchain"
)

// MemoizingLoader ...
type MemoizingLoader struct {
	loader         blockchain.Loader
	loadingResults LRUCache
}

// NewMemoizingLoader ...
func NewMemoizingLoader(
	maximalCacheSize int,
	loader blockchain.Loader,
) MemoizingLoader {
	return MemoizingLoader{
		loader:         loader,
		loadingResults: NewLRUCache(maximalCacheSize),
	}
}

// LoadBlocks ...
func (loader MemoizingLoader) LoadBlocks(cursor interface{}, count int) (
	blocks blockchain.BlockGroup,
	nextCursor interface{},
	err error,
) {
	parameters := Parameters{Cursor: cursor, Count: count}
	results, isFound := loader.loadingResults.Get(parameters)
	if isFound {
		return results.Blocks, results.NextCursor, nil
	}

	blocks, nextCursor, err = loader.loader.LoadBlocks(cursor, count)
	if err != nil {
		return nil, nil, err
	}

	results = Results{Blocks: blocks, NextCursor: nextCursor}
	loader.loadingResults.Set(parameters, results)

	return blocks, nextCursor, nil
}
