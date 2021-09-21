package blockchain

import (
	"sync"
)

type loadingParameters struct {
	cursor interface{}
	count  int
}

type loadingResult struct {
	blocks     BlockGroup
	nextCursor interface{}
}

// MemoizingLoader ...
type MemoizingLoader struct {
	loader         Loader
	loadingResults *sync.Map
}

// NewMemoizingLoader ...
func NewMemoizingLoader(loader Loader) MemoizingLoader {
	return MemoizingLoader{
		loader:         loader,
		loadingResults: new(sync.Map),
	}
}
