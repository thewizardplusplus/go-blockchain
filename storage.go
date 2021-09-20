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

//go:generate mockery --name=Loader --inpackage --case=underscore --testonly

// Loader ...
type Loader interface {
	LoadBlocks(cursor interface{}, count int) (
		blocks BlockGroup,
		nextCursor interface{},
		err error,
	)
}

// LoadStorage ...
func LoadStorage(
	storage Storage,
	loader Loader,
	initialCursor interface{},
	chunkSize int,
) (lastCursor interface{}, err error) {
	cursor := initialCursor
	for {
		blocks, nextCursor, err := loader.LoadBlocks(cursor, chunkSize)
		if err != nil {
			const message = "unable to load the blocks corresponding to cursor %v"
			return cursor, errors.Wrapf(err, message, cursor)
		}
		if len(blocks) == 0 {
			break
		}

		for index, block := range blocks {
			if err := storage.StoreBlock(block); err != nil {
				const message = "unable to store block #%d " +
					"from the blocks corresponding to cursor %v"
				return cursor, errors.Wrapf(err, message, index, cursor)
			}
		}

		cursor = nextCursor
	}

	return cursor, nil
}
