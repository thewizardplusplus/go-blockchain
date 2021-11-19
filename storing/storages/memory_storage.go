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
	if !storage.isSorted {
		sort.Slice(storage.blocks, func(i int, j int) bool {
			// descending order
			return storage.blocks[j].Timestamp.After(storage.blocks[i].Timestamp)
		})

		storage.isSorted = true
	}

	loader := loaders.MemoryLoader(storage.blocks)
	return loader.LoadBlocks(cursor, count)
}

// LoadLastBlock ...
func (storage MemoryStorage) LoadLastBlock() (blockchain.Block, error) {
	if storage.isEmpty() {
		return blockchain.Block{}, blockchain.ErrEmptyStorage
	}

	return storage.lastBlock, nil
}

// StoreBlock ...
func (storage *MemoryStorage) StoreBlock(block blockchain.Block) error {
	// this check should follow before appending the new block
	if storage.isEmpty() || block.Timestamp.After(storage.lastBlock.Timestamp) {
		storage.lastBlock = block
	}

	storage.blocks = append(storage.blocks, block)
	storage.isSorted = false

	return nil
}

func (storage MemoryStorage) isEmpty() bool {
	return len(storage.blocks) == 0
}
