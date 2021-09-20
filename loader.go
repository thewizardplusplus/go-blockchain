package blockchain

import (
	"github.com/pkg/errors"
)

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
	storage GroupStorage,
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

		if err := storage.StoreBlockGroup(blocks); err != nil {
			const message = "unable to store the blocks corresponding to cursor %v"
			return cursor, errors.Wrapf(err, message, cursor)
		}

		cursor = nextCursor
	}

	return cursor, nil
}
