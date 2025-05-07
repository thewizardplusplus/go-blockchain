package blockchain

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/mo"
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
//
// Deprecated: Use [NewBlockchainEx] instead.
func NewBlockchain(
	genesisBlockData Data,
	dependencies Dependencies,
) (*Blockchain, error) {
	blockchain, err := NewBlockchainEx(context.Background(), NewBlockchainExParams{
		Dependencies:     dependencies,
		GenesisBlockData: mo.EmptyableToOption(genesisBlockData),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create a new blockchain: %w", err)
	}

	return blockchain, nil
}

// NewBlockchainExParams ...
type NewBlockchainExParams struct {
	Dependencies     Dependencies
	GenesisBlockData mo.Option[Data]
}

// NewBlockchainEx ...
func NewBlockchainEx(
	ctx context.Context,
	params NewBlockchainExParams,
) (*Blockchain, error) {
	lastBlock, err := params.Dependencies.Storage.LoadLastBlock()
	if err != nil &&
		(!errors.Is(err, ErrEmptyStorage) || params.GenesisBlockData.IsAbsent()) {
		return nil, fmt.Errorf("unable to load the last block: %w", err)
	}

	if errors.Is(err, ErrEmptyStorage) && params.GenesisBlockData.IsPresent() {
		genesisBlock, err := NewGenesisBlockEx(ctx, NewGenesisBlockExParams{
			Dependencies: params.Dependencies.BlockDependencies,
			Data:         params.GenesisBlockData.MustGet(),
		})
		if err != nil {
			return nil, fmt.Errorf("unable to create a new genesis block: %w", err)
		}

		if err = params.Dependencies.Storage.StoreBlock(genesisBlock); err != nil {
			return nil, fmt.Errorf("unable to store the genesis block: %w", err)
		}

		lastBlock = genesisBlock
	}

	blockchain := &Blockchain{
		dependencies: params.Dependencies,
		lastBlock:    lastBlock,
	}
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
//
// Deprecated: Use [AddBlockEx] instead.
func (blockchain *Blockchain) AddBlock(data Data) error {
	if err := blockchain.AddBlockEx(context.Background(), data); err != nil {
		return fmt.Errorf("unable to add a new block: %w", err)
	}

	return nil
}

// AddBlockEx ...
func (blockchain *Blockchain) AddBlockEx(ctx context.Context, data Data) error {
	block, err := NewBlockEx(ctx, NewBlockExParams{
		Dependencies: blockchain.dependencies.BlockDependencies,
		Data:         data,
		PrevBlock:    mo.Some(blockchain.lastBlock),
	})
	if err != nil {
		return fmt.Errorf("unable to create a new block: %w", err)
	}

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
