package blockchain

import (
	"time"

	"github.com/pkg/errors"
)

// ValidationMode ...
type ValidationMode int

// ...
const (
	AsFullBlockchain ValidationMode = iota
	AsBlockchainChunk
)

// BlockGroup ...
type BlockGroup []Block

// IsValid ...
func (blocks BlockGroup) IsValid(
	prependedChunk BlockGroup,
	validationMode ValidationMode,
	proofer Proofer,
) error {
	if len(blocks) == 0 {
		return nil
	}

	if len(prependedChunk) != 0 {
		prevBlock := &blocks[0]
		err := prependedChunk.IsLastBlockValid(prevBlock, AsBlockchainChunk, proofer)
		if err != nil {
			return errors.Wrap(err, "the prepended chunk is not valid")
		}
	}

	for index, block := range blocks[:len(blocks)-1] {
		prevBlock := &blocks[index+1]
		if err := block.IsValid(prevBlock, proofer); err != nil {
			return errors.Wrapf(err, "block #%d is not valid", index)
		}
	}

	if err := blocks.IsLastBlockValid(nil, validationMode, proofer); err != nil {
		return errors.Wrap(err, "the last block is not valid")
	}

	return nil
}

// IsLastBlockValid ...
func (blocks BlockGroup) IsLastBlockValid(
	prevBlock *Block,
	validationMode ValidationMode,
	proofer Proofer,
) error {
	var err error
	switch lastBlock := blocks[len(blocks)-1]; validationMode {
	case AsFullBlockchain:
		if err = lastBlock.IsValidGenesisBlock(proofer); err != nil {
			err = errors.Wrap(err, "the last block was validated as a genesis block")
		}
	case AsBlockchainChunk:
		err = lastBlock.IsValid(prevBlock, proofer)
	}

	return err
}

// FindDifferences ...
func (blocks BlockGroup) FindDifferences(anotherBlocks BlockGroup) (
	leftIndex int,
	rightIndex int,
	hasMatch bool,
) {
	timestampIndexMap := make(map[time.Time]int)
	for index, block := range anotherBlocks {
		timestampIndexMap[normalizeTimestamp(block.Timestamp)] = index
	}

	for index, block := range blocks {
		anotherIndex, isTimestampFound :=
			timestampIndexMap[normalizeTimestamp(block.Timestamp)]
		if isTimestampFound && anotherBlocks[anotherIndex].IsEqual(block) == nil {
			return index, anotherIndex, true
		}
	}

	return 0, 0, false
}

// Difficulty ...
func (blocks BlockGroup) Difficulty(proofer Proofer) (int, error) {
	var totalDifficulty int
	for index, block := range blocks {
		difficulty, err := proofer.Difficulty(block.Hash)
		if err != nil {
			return 0, errors.Wrapf(
				err,
				"unable to calculate the difficulty of the block #%d",
				index,
			)
		}

		totalDifficulty += difficulty
	}

	return totalDifficulty, nil
}

func normalizeTimestamp(timestamp time.Time) time.Time {
	return timestamp.
		In(time.UTC). // set the same location for all timestamps
		Truncate(0)   // strip monotonic clock reading
}
