package blockchain

import (
	"time"
)

//go:generate mockery --name=Hasher --inpackage --case=underscore --testonly

// Hasher ...
type Hasher interface {
	Hash() string
}

// Clock ...
type Clock func() time.Time

//go:generate mockery --name=Proofer --inpackage --case=underscore --testonly

// Proofer ...
type Proofer interface {
	Hash(block Block) string
	Validate(block Block) bool
}

// Dependencies ...
type Dependencies struct {
	Clock   Clock
	Proofer Proofer
}

// Block ...
type Block struct {
	Timestamp time.Time
	Data      Hasher
	Hash      string
	PrevHash  string
}

// NewBlock ...
func NewBlock(data Hasher, prevBlock Block, dependencies Dependencies) Block {
	block := Block{
		Timestamp: dependencies.Clock(),
		Data:      data,
		PrevHash:  prevBlock.Hash,
	}
	block.Hash = dependencies.Proofer.Hash(block)

	return block
}
