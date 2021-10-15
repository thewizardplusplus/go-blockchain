package blockchain

import (
	"github.com/pkg/errors"
)

// ErrEmptyStorage ...
var ErrEmptyStorage = errors.New("empty storage")

//go:generate mockery --name=Storage --inpackage --case=underscore --testonly

// Storage ...
type Storage interface {
	LoadLastBlock() (Block, error)
	StoreBlock(block Block) error
}

//go:generate mockery --name=GroupStorage --inpackage --case=underscore --testonly

// GroupStorage ...
type GroupStorage interface {
	Storage

	StoreBlockGroup(blocks BlockGroup) error
}
