package blockchain

import (
	"github.com/pkg/errors"
)

// ErrEmptyStorage ...
var ErrEmptyStorage = errors.New("empty storage")

// Storage ...
type Storage interface {
	LoadLastBlock() (Block, error)
	StoreBlock(block Block) error
}

// Dependencies ...
type Dependencies struct {
	BlockDependencies

	Storage Storage
}

// Blockchain ...
type Blockchain struct {
	dependencies Dependencies
	lastBlock    Block
}
