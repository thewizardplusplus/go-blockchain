package storages

import (
	"container/heap"

	"github.com/thewizardplusplus/go-blockchain"
)

// MemoryStorage ...
type MemoryStorage struct {
	blocks blockchain.BlockGroup
}

// Blocks ...
func (storage MemoryStorage) Blocks() blockchain.BlockGroup {
	blocksCopy := make(blockchain.BlockGroup, len(storage.blocks))
	copy(blocksCopy, storage.blocks)
	heap.Init((*BlockPriorityQueue)(&blocksCopy))

	blocks := make(blockchain.BlockGroup, len(storage.blocks))
	for len(blocksCopy) != 0 {
		targetIndex := len(blocksCopy) - 1
		lastBlock := heap.Pop((*BlockPriorityQueue)(&blocksCopy))
		blocks[targetIndex] = lastBlock.(blockchain.Block)
	}

	return blocks
}

// LoadLastBlock ...
func (storage *MemoryStorage) LoadLastBlock() (blockchain.Block, error) {
	if len(storage.blocks) == 0 {
		return blockchain.Block{}, blockchain.ErrEmptyStorage
	}

	lastBlock := heap.Pop((*BlockPriorityQueue)(&storage.blocks))
	// restore the popped block
	heap.Push((*BlockPriorityQueue)(&storage.blocks), lastBlock)

	return lastBlock.(blockchain.Block), nil
}

// StoreBlock ...
func (storage *MemoryStorage) StoreBlock(block blockchain.Block) error {
	heap.Push((*BlockPriorityQueue)(&storage.blocks), block)
	return nil
}
