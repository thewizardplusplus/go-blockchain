package blockchain

import (
	"time"
)

// Hasher ...
type Hasher interface {
	Hash() string
}

// Block ...
type Block struct {
	Timestamp time.Time
	Data      Hasher
	Hash      string
	PrevHash  string
}
