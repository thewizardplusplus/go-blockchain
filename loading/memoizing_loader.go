package loading

import (
	"sync"

	"github.com/thewizardplusplus/go-blockchain"
)

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
	parameters := Parameters{Cursor: cursor, Count: count}
	results, ok := loader.loadingResults.Load(parameters)
	if ok {
		return results.(Results).Blocks, results.(Results).NextCursor, nil
	}

	blocks, nextCursor, err = loader.loader.LoadBlocks(cursor, count)
	if err != nil {
		return nil, nil, err
	}

	results = Results{Blocks: blocks, NextCursor: nextCursor}
	loader.loadingResults.Store(parameters, results)

	return blocks, nextCursor, nil
}
