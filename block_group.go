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
