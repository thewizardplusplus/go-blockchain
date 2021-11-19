package blockchain

import (
	"fmt"

	"github.com/pkg/errors"
)

// Dependencies ...
type Dependencies struct {
	BlockDependencies

	Storage GroupStorage
}

// Blockchain ...
type Blockchain struct {
	dependencies Dependencies
	lastBlock    Block
}

// NewBlockchain ...
func NewBlockchain(
	genesisBlockData fmt.Stringer,
	dependencies Dependencies,
) (*Blockchain, error) {
	lastBlock, err := dependencies.Storage.LoadLastBlock()
	switch {
	case err == nil:
	case errors.Cause(err) == ErrEmptyStorage && genesisBlockData != nil:
		genesisBlock :=
			NewGenesisBlock(genesisBlockData, dependencies.BlockDependencies)
		if err = dependencies.Storage.StoreBlock(genesisBlock); err != nil {
			return nil, errors.Wrap(err, "unable to store the genesis block")
		}

		lastBlock = genesisBlock
	default:
		return nil, errors.Wrap(err, "unable to load the last block")
	}

	blockchain := &Blockchain{dependencies: dependencies, lastBlock: lastBlock}
	return blockchain, nil
}

// LoadBlocks ...
func (blockchain Blockchain) LoadBlocks(cursor interface{}, count int) (
	blocks BlockGroup,
	nextCursor interface{},
	err error,
) {
	return blockchain.dependencies.Storage.LoadBlocks(cursor, count)
}

// AddBlock ...
func (blockchain *Blockchain) AddBlock(data fmt.Stringer) error {
	block := NewBlock(
		data,
		blockchain.lastBlock,
		blockchain.dependencies.BlockDependencies,
	)
	if err := blockchain.dependencies.Storage.StoreBlock(block); err != nil {
		return errors.Wrap(err, "unable to store the block")
	}

	blockchain.lastBlock = block
	return nil
}
