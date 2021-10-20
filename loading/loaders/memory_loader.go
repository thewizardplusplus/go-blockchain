package loaders

import (
	"github.com/thewizardplusplus/go-blockchain"
)

// MemoryLoader ...
type MemoryLoader blockchain.BlockGroup

// LoadBlocks ...
func (loader MemoryLoader) LoadBlocks(cursor interface{}, count int) (
	blocks blockchain.BlockGroup,
	nextCursor interface{},
	err error,
) {
	blocks = blockchain.BlockGroup(loader)

	var startIndex int
	if cursor != nil {
		startIndex = cursor.(int)
	}

	endIndex := startIndex + count
	if maximalEndIndex := len(blocks); endIndex > maximalEndIndex {
		endIndex = maximalEndIndex
	}

	if maximalStartIndex := len(blocks) - 1; startIndex > maximalStartIndex {
		return blockchain.BlockGroup{}, endIndex, nil
	}

	blocks = blocks[startIndex:endIndex]
	return blocks, endIndex, nil
}
