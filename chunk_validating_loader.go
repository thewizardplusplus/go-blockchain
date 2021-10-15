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

	err = blocks.IsValid(nil, AsBlockchainChunk, loader.Dependencies)
	if err != nil {
		const message = "the blocks corresponding to cursor %v are not valid"
		return nil, nil, errors.Wrapf(err, message, cursor)
	}

	return blocks, nextCursor, nil
}
