package storages

import (
	"github.com/thewizardplusplus/go-blockchain"
)

// MemoryStorage ...
type MemoryStorage struct {
	blocks blockchain.BlockGroup
}

// Blocks ...
func (storage MemoryStorage) Blocks() blockchain.BlockGroup {
	return storage.blocks
}

// LoadLastBlock ...
func (storage MemoryStorage) LoadLastBlock() (blockchain.Block, error) {
	if len(storage.blocks) == 0 {
		return blockchain.Block{}, blockchain.ErrEmptyStorage
	}

	lastBlock := storage.blocks[len(storage.blocks)-1]
	return lastBlock, nil
}

// StoreBlock ...
func (storage *MemoryStorage) StoreBlock(block blockchain.Block) error {
	storage.blocks = append(storage.blocks, block)
	return nil
}
