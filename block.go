package blockchain

import (
	"time"
)

// Hasher ...
type Hasher interface {
	Hash() string
}

// Proofer ...
type Proofer interface {
	Hash(block Block) string
	Validate(block Block) bool
}

// Block ...
type Block struct {
	Timestamp time.Time
	Data      Hasher
	Hash      string
	PrevHash  string
}
