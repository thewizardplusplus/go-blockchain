package blockchain

import (
	"errors"
	"fmt"
)

// ErrEqualDifficulties ...
var ErrEqualDifficulties = errors.New("equal difficulties")

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
	genesisBlockData Data,
	dependencies Dependencies,
) (*Blockchain, error) {
	lastBlock, err := dependencies.Storage.LoadLastBlock()
	switch {
	case err == nil:
	case errors.Is(err, ErrEmptyStorage) && genesisBlockData != nil:
		genesisBlock :=
			NewGenesisBlock(genesisBlockData, dependencies.BlockDependencies)
		if err = dependencies.Storage.StoreBlock(genesisBlock); err != nil {
			return nil, fmt.Errorf("unable to store the genesis block: %w", err)
		}

		lastBlock = genesisBlock
	default:
		return nil, fmt.Errorf("unable to load the last block: %w", err)
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
func (blockchain *Blockchain) AddBlock(data Data) error {
	block := NewBlock(
		data,
		blockchain.lastBlock,
		blockchain.dependencies.BlockDependencies,
	)
	if err := blockchain.dependencies.Storage.StoreBlock(block); err != nil {
		return fmt.Errorf("unable to store the block: %w", err)
	}

	blockchain.lastBlock = block
	return nil
}

// Merge ...
func (blockchain *Blockchain) Merge(loader Loader, chunkSize int) error {
	leftDifferences, rightDifferences, err :=
		FindDifferences(blockchain, loader, chunkSize)
	if err != nil {
		return fmt.Errorf("unable to find differences: %w", err)
	}

	leftDifficulty, err :=
		leftDifferences.Difficulty(blockchain.dependencies.Proofer)
	if err != nil {
		return fmt.Errorf(
			"unable to calculate the difficulty of the left differences: %w",
			err,
		)
	}

	rightDifficulty, err :=
		rightDifferences.Difficulty(blockchain.dependencies.Proofer)
	if err != nil {
		return fmt.Errorf(
			"unable to calculate the difficulty of the right differences: %w",
			err,
		)
	}

	if leftDifficulty > rightDifficulty {
		return nil
	}
	if leftDifficulty == rightDifficulty {
		return ErrEqualDifficulties
	}

	// if leftDifficulty < rightDifficulty...
	if err = blockchain.dependencies.Storage.
		DeleteBlockGroup(leftDifferences); err != nil {
		return fmt.Errorf("unable to delete the left differences: %w", err)
	}

	if err = blockchain.dependencies.Storage.
		StoreBlockGroup(rightDifferences); err != nil {
		return fmt.Errorf("unable to store the right differences: %w", err)
	}

	lastBlock, err := blockchain.dependencies.Storage.LoadLastBlock()
	if err != nil {
		return fmt.Errorf("unable to load the last block: %w", err)
	}
	blockchain.lastBlock = lastBlock

	return nil
}
