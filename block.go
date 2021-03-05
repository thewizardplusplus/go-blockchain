package blockchain

import (
	"time"
)

// Hasher ...
type Hasher interface {
	Hash() string
}

// Clock ...
type Clock func() time.Time

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
