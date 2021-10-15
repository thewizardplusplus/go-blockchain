package loading

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-blockchain"
)

// LoadStorage ...
func LoadStorage(
	storage blockchain.GroupStorage,
	loader blockchain.Loader,
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
