package blockchain

import (
	"github.com/pkg/errors"
)

// ChunkValidatingLoader ...
type ChunkValidatingLoader struct {
	Loader       Loader
	Dependencies BlockDependencies
}

// LoadBlocks ...
func (loader ChunkValidatingLoader) LoadBlocks(cursor interface{}, count int) (
	blocks BlockGroup,
	nextCursor interface{},
	err error,
) {
	blocks, nextCursor, err = loader.Loader.LoadBlocks(cursor, count)
	if err != nil {
		return nil, nil, err
	}

	if !blocks.IsValid(nil, AsBlockchainChunk, loader.Dependencies) {
		const message = "the blocks corresponding to cursor %v are not valid"
		return nil, nil, errors.Errorf(message, cursor)
	}

	return blocks, nextCursor, nil
}
