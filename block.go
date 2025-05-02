package blockchain

import (
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
func NewBlock(
	data Data,
	prevBlock Block,
	dependencies BlockDependencies,
) Block {
	block := Block{
		Timestamp: dependencies.Clock(),
		Data:      data,
		PrevHash:  prevBlock.Hash,
	}
	block.Hash = dependencies.Proofer.Hash(block)

	return block
}

// NewGenesisBlock ...
func NewGenesisBlock(data Data, dependencies BlockDependencies) Block {
	return NewBlock(data, Block{}, dependencies)
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
