package storages

import (
	"github.com/thewizardplusplus/go-blockchain"
)

// MemoryStorage ...
type MemoryStorage struct {
	blocks    blockchain.BlockGroup
	lastBlock blockchain.Block
}

// Blocks ...
func (storage MemoryStorage) Blocks() blockchain.BlockGroup {
	return storage.blocks
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
	return nil
}

func (storage MemoryStorage) isEmpty() bool {
	return len(storage.blocks) == 0
}
