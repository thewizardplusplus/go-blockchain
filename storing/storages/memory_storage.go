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

// NewMemoryStorage ...
func NewMemoryStorage(blocks blockchain.BlockGroup) *MemoryStorage {
	storage := &MemoryStorage{blocks: blocks}
	for _, block := range storage.blocks {
		if block.Timestamp.After(storage.lastBlock.Timestamp) {
			storage.lastBlock = block
		}
	}

	return storage
}

// LoadBlocks ...
func (storage *MemoryStorage) LoadBlocks(cursor interface{}, count int) (
	blocks blockchain.BlockGroup,
	nextCursor interface{},
	err error,
) {
	storage.sortIfNeed()

	loader := loaders.MemoryLoader(storage.blocks)
	// the memory loader never returns an error
	blocks, nextCursor, _ = loader.LoadBlocks(cursor, count) // nolint: gosec

	copiedBlocks := make(blockchain.BlockGroup, len(blocks))
	copy(copiedBlocks, blocks)

	return copiedBlocks, nextCursor, nil
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

	index := sort.Search(len(storage.blocks), func(index int) bool {
		return !storage.blocks[index].Timestamp.
			After(block.Timestamp) // before or equal
	})
	if index == len(storage.blocks) ||
		storage.blocks[index].IsEqual(block) != nil {
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
