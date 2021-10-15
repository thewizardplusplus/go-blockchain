package blockchain

import (
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
	dependencies BlockDependencies,
) error {
	if len(blocks) == 0 {
		return nil
	}

	if len(prependedChunk) != 0 {
		prevBlock := &blocks[0]
		err := prependedChunk.IsLastBlockValid(
			prevBlock,
			AsBlockchainChunk,
			dependencies,
		)
		if err != nil {
			return errors.Wrap(err, "the prepended chunk is not valid")
		}
	}

	for index, block := range blocks[:len(blocks)-1] {
		prevBlock := &blocks[index+1]
		if err := block.IsValid(prevBlock, dependencies); err != nil {
			return errors.Wrapf(err, "block #%d is not valid", index)
		}
	}

	err := blocks.IsLastBlockValid(nil, validationMode, dependencies)
	if err != nil {
		return errors.Wrap(err, "the last block is not valid")
	}

	return nil
}

// IsLastBlockValid ...
func (blocks BlockGroup) IsLastBlockValid(
	prevBlock *Block,
	validationMode ValidationMode,
	dependencies BlockDependencies,
) error {
	var err error
	lastBlock := blocks[len(blocks)-1]
	switch validationMode {
	case AsFullBlockchain:
		err = lastBlock.IsValidGenesisBlock(dependencies)
		if err != nil {
			err = errors.Wrap(err, "the last block was validated as a genesis block")
		}
	case AsBlockchainChunk:
		err = lastBlock.IsValid(prevBlock, dependencies)
	}

	return err
}
