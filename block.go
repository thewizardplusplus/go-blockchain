package blockchain

import (
	"fmt"
	"time"
)

// Clock ...
type Clock func() time.Time

//go:generate mockery --name=Proofer --inpackage --case=underscore --testonly

// Proofer ...
type Proofer interface {
	Hash(block Block) string
	Validate(block Block) bool
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
) bool {
	var prevTimestamp time.Time
	if prevBlock != nil {
		prevTimestamp = prevBlock.Timestamp
	}
	if !block.Timestamp.After(prevTimestamp) {
		return false
	}

	if prevBlock != nil && block.PrevHash != prevBlock.Hash {
		return false
	}

	return dependencies.Proofer.Validate(block)
}

// IsValidGenesisBlock ...
func (block Block) IsValidGenesisBlock(dependencies BlockDependencies) bool {
	return block.IsValid(&Block{}, dependencies)
}
