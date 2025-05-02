package loading

import (
	"fmt"

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
			const message = "unable to load the blocks corresponding to cursor %v: %w"
			return cursor, fmt.Errorf(message, cursor, err)
		}
		if len(blocks) == 0 {
			break
		}

		if err := storage.StoreBlockGroup(blocks); err != nil {
			const message = "unable to store the blocks corresponding to cursor %v: %w"
			return cursor, fmt.Errorf(message, cursor, err)
		}

		cursor = nextCursor
	}

	return cursor, nil
}
