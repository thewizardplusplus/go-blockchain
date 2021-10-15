package loading

import (
	"sync"

	"github.com/thewizardplusplus/go-blockchain"
)

type loadingParameters struct {
	cursor interface{}
	count  int
}

type loadingResult struct {
	blocks     blockchain.BlockGroup
	nextCursor interface{}
}

// MemoizingLoader ...
type MemoizingLoader struct {
	loader         blockchain.Loader
	loadingResults *sync.Map
}

// NewMemoizingLoader ...
func NewMemoizingLoader(loader blockchain.Loader) MemoizingLoader {
	return MemoizingLoader{
		loader:         loader,
		loadingResults: new(sync.Map),
	}
}

// LoadBlocks ...
func (loader MemoizingLoader) LoadBlocks(cursor interface{}, count int) (
	blocks blockchain.BlockGroup,
	nextCursor interface{},
	err error,
) {
	parameters := loadingParameters{cursor: cursor, count: count}
	results, ok := loader.loadingResults.Load(parameters)
	if ok {
		return results.(loadingResult).blocks, results.(loadingResult).nextCursor, nil
	}

	blocks, nextCursor, err = loader.loader.LoadBlocks(cursor, count)
	if err != nil {
		return nil, nil, err
	}

	results = loadingResult{blocks: blocks, nextCursor: nextCursor}
	loader.loadingResults.Store(parameters, results)

	return blocks, nextCursor, nil
}
