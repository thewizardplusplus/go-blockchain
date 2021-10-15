package blockchain

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Clock ...
type Clock func() time.Time

//go:generate mockery --name=Proofer --inpackage --case=underscore --testonly

// Proofer ...
type Proofer interface {
	Hash(block Block) string
	Validate(block Block) error
}

// BlockDependencies ...
type BlockDependencies struct {
	Clock   Clock
	Proofer Proofer
}

// Block ...
type Block struct {
	Timestamp time.Time
	Data      fmt.Stringer
	Hash      string
	PrevHash  string
}

// NewBlock ...
func NewBlock(
	data fmt.Stringer,
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
func NewGenesisBlock(data fmt.Stringer, dependencies BlockDependencies) Block {
	return NewBlock(data, Block{}, dependencies)
}

// MergedData ...
func (block Block) MergedData() string {
	return block.Timestamp.String() + block.Data.String() + block.PrevHash
}

// IsValid ...
func (block Block) IsValid(
	prevBlock *Block,
	dependencies BlockDependencies,
) error {
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

	if err := dependencies.Proofer.Validate(block); err != nil {
		return errors.Wrap(err, "the validation via the proofer was failed")
	}

	return nil
}

// IsValidGenesisBlock ...
func (block Block) IsValidGenesisBlock(dependencies BlockDependencies) error {
	return block.IsValid(&Block{}, dependencies)
}
