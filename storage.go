package blockchain

import (
	"github.com/pkg/errors"
)

// ErrEmptyStorage ...
var ErrEmptyStorage = errors.New("empty storage")

// Storage ...
type Storage interface {
	Loader

	LoadLastBlock() (Block, error)
	StoreBlock(block Block) error
	DeleteBlock(block Block) error
}

//go:generate mockery --name=GroupStorage --inpackage --case=underscore --testonly

// GroupStorage ...
type GroupStorage interface {
	Storage

	StoreBlockGroup(blocks BlockGroup) error
	DeleteBlockGroup(blocks BlockGroup) error
}
