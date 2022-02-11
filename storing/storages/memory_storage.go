package storages

import (
	"sort"

	"github.com/thewizardplusplus/go-blockchain"
	"github.com/thewizardplusplus/go-blockchain/loading/loaders"
)

// MemoryStorage ...
type MemoryStorage struct {
	blocks    blockchain.BlockGroup
	lastBlock blockchain.Block
	isSorted  bool
}

// LoadBlocks ...
func (storage *MemoryStorage) LoadBlocks(cursor interface{}, count int) (
	blocks blockchain.BlockGroup,
	nextCursor interface{},
	err error,
) {
	storage.sortIfNeed()

	loader := loaders.MemoryLoader(storage.blocks)
	return loader.LoadBlocks(cursor, count)
}

// LoadLastBlock ...
func (storage MemoryStorage) LoadLastBlock() (blockchain.Block, error) {
	if len(storage.blocks) == 0 {
		return blockchain.Block{}, blockchain.ErrEmptyStorage
	}

	return storage.lastBlock, nil
}

// StoreBlock ...
func (storage *MemoryStorage) StoreBlock(block blockchain.Block) error {
	// this check should follow before appending the new block
	if len(storage.blocks) == 0 ||
		block.Timestamp.After(storage.lastBlock.Timestamp) {
		storage.lastBlock = block
	}

	storage.blocks = append(storage.blocks, block)
	storage.isSorted = false

	return nil
}

// DeleteBlock ...
func (storage *MemoryStorage) DeleteBlock(block blockchain.Block) error {
	storage.sortIfNeed()

	index, isFound := storage.blocks.FindBlock(block)
	if !isFound {
		return nil
	}

	// https://github.com/golang/go/wiki/SliceTricks#delete
	copiedCount := copy(storage.blocks[index:], storage.blocks[index+1:])
	storage.blocks = storage.blocks[:index+copiedCount]

	if len(storage.blocks) != 0 {
		storage.lastBlock = storage.blocks[0]
	}

	return nil
}

func (storage *MemoryStorage) sortIfNeed() {
	if storage.isSorted {
		return
	}

	sort.Slice(storage.blocks, func(i int, j int) bool {
		return storage.blocks[i].Timestamp.
			After(storage.blocks[j].Timestamp) // descending order
	})
	storage.isSorted = true
}
