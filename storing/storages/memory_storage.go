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

// DeleteBlock ...
func (storage *MemoryStorage) DeleteBlock(block blockchain.Block) error {
	storage.sortIfNeed()

	index := sort.Search(len(storage.blocks), func(index int) bool {
		// after or equal
		return !storage.blocks[index].Timestamp.Before(block.Timestamp)
	})
	if index < len(storage.blocks) &&
		storage.blocks[index].Timestamp.Equal(block.Timestamp) {
		// https://github.com/golang/go/wiki/SliceTricks#delete
		copiedCount := copy(storage.blocks[index:], storage.blocks[index+1:])
		storage.blocks = storage.blocks[:index+copiedCount]

		if !storage.isEmpty() {
			storage.lastBlock = storage.blocks[len(storage.blocks)-1]
		}
	}

	return nil
}

func (storage MemoryStorage) isEmpty() bool {
	return len(storage.blocks) == 0
}

func (storage *MemoryStorage) sortIfNeed() {
	if storage.isSorted {
		return
	}

	sort.Slice(storage.blocks, func(i int, j int) bool {
		// descending order
		return storage.blocks[j].Timestamp.After(storage.blocks[i].Timestamp)
	})
	storage.isSorted = true
}
