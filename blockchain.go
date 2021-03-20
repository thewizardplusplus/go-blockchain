package blockchain

import (
	"github.com/pkg/errors"
)

// ErrEmptyStorage ...
var ErrEmptyStorage = errors.New("empty storage")

//go:generate mockery --name=Storage --inpackage --case=underscore --testonly

// Storage ...
type Storage interface {
	LoadLastBlock() (Block, error)
	StoreBlock(block Block) error
}

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
		if err := dependencies.Storage.StoreBlock(genesisBlock); err != nil {
			return nil, errors.Wrap(err, "unable to store the genesis block")
		}

		lastBlock = genesisBlock
	default:
		return nil, errors.Wrap(err, "unable to load the last block")
	}

	blockchain := &Blockchain{dependencies: dependencies, lastBlock: lastBlock}
	return blockchain, nil
}
