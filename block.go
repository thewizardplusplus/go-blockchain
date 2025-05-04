package blockchain

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Clock ...
type Clock func() time.Time

//go:generate mockery --name=Proofer --inpackage --case=underscore --testonly

// Proofer ...
type Proofer interface {
	Hash(block Block) string
	HashEx(ctx context.Context, block Block) (string, error)
	Validate(block Block) error
	Difficulty(hash string) (int, error)
}

// BlockDependencies ...
type BlockDependencies struct {
	Clock   Clock
	Proofer Proofer
}

// Block ...
type Block struct {
	Timestamp time.Time
	Data      Data
	Hash      string
	PrevHash  string
}

// NewBlock ...
//
// Deprecated: Use [NewBlockEx] instead.
func NewBlock(
	data Data,
	prevBlock Block,
	dependencies BlockDependencies,
) Block {
	block, _ := NewBlockEx(context.Background(), NewBlockExParams{
		Clock:     dependencies.Clock,
		Data:      data,
		PrevBlock: prevBlock,
		Proofer:   dependencies.Proofer,
	})
	return block
}

// NewBlockExParams ...
type NewBlockExParams struct {
	Clock     Clock
	Data      Data
	PrevBlock Block
	Proofer   Proofer
}

// NewBlockEx ...
func NewBlockEx(ctx context.Context, params NewBlockExParams) (Block, error) {
	block := Block{
		Timestamp: params.Clock(),
		Data:      params.Data,
		PrevHash:  params.PrevBlock.Hash,
	}

	var err error
	block.Hash, err = params.Proofer.HashEx(ctx, block)
	if err != nil {
		return Block{}, fmt.Errorf("unable to hash a new block: %w", err)
	}

	return block, nil
}

// NewGenesisBlock ...
//
// Deprecated: Use [NewGenesisBlockEx] instead.
func NewGenesisBlock(data Data, dependencies BlockDependencies) Block {
	genesisBlock, _ := NewGenesisBlockEx(
		context.Background(),
		NewGenesisBlockExParams{
			Clock:   dependencies.Clock,
			Data:    data,
			Proofer: dependencies.Proofer,
		},
	)
	return genesisBlock
}

// NewGenesisBlockExParams ...
type NewGenesisBlockExParams struct {
	Clock   Clock
	Data    Data
	Proofer Proofer
}

// NewGenesisBlockEx ...
func NewGenesisBlockEx(
	ctx context.Context,
	params NewGenesisBlockExParams,
) (Block, error) {
	genesisBlock, err := NewBlockEx(ctx, NewBlockExParams{
		Clock:     params.Clock,
		Data:      params.Data,
		PrevBlock: Block{},
		Proofer:   params.Proofer,
	})
	if err != nil {
		return Block{}, fmt.Errorf("unable to create a new block: %w", err)
	}

	return genesisBlock, nil
}

// MergedData ...
func (block Block) MergedData() string {
	return block.Timestamp.String() + block.Data.String() + block.PrevHash
}

// IsEqual ...
func (block Block) IsEqual(anotherBlock Block) error {
	if !block.Timestamp.Equal(anotherBlock.Timestamp) {
		return errors.New("timestamps are not equal")
	}
	if !block.Data.Equal(anotherBlock.Data) {
		return errors.New("data are not equal")
	}
	if block.Hash != anotherBlock.Hash {
		return errors.New("hashes are not equal")
	}
	if block.PrevHash != anotherBlock.PrevHash {
		return errors.New("previous hashes are not equal")
	}
	return nil
}

// IsValid ...
func (block Block) IsValid(prevBlock *Block, proofer Proofer) error {
	var prevTimestamp time.Time
	if prevBlock != nil {
		prevTimestamp = prevBlock.Timestamp
	}
	if !block.Timestamp.After(prevTimestamp) {
		return errors.New("the timestamp is not greater than the previous one")
	}

	if prevBlock != nil && block.PrevHash != prevBlock.Hash {
		return errors.New(
			"the previous hash is not equal to the hash of the previous block",
		)
	}

	if err := proofer.Validate(block); err != nil {
		return fmt.Errorf("the validation via the proofer was failed: %w", err)
	}

	return nil
}

// IsValidGenesisBlock ...
func (block Block) IsValidGenesisBlock(proofer Proofer) error {
	return block.IsValid(&Block{}, proofer)
}
