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

// Loader ...
type Loader interface {
	LoadBlocks(cursor interface{}, count int) (
		blocks BlockGroup,
		nextCursor interface{},
		err error,
	)
}
