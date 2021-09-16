package blockchain

import (
	"github.com/pkg/errors"
)

// Dependencies ...
type Dependencies struct {
	BlockDependencies

	Storage Storage
}

// Blockchain ...
type Blockchain struct {
	dependencies Dependencies
	lastBlock    Block
}

// NewBlockchain ...
func NewBlockchain(
	genesisBlockData Hasher,
	dependencies Dependencies,
) (*Blockchain, error) {
	lastBlock, err := dependencies.Storage.LoadLastBlock()
	switch {
	case err == nil:
	case errors.Cause(err) == ErrEmptyStorage:
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

// AddBlock ...
func (blockchain *Blockchain) AddBlock(data Hasher) error {
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
