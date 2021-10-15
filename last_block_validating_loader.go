package blockchain

import (
	"github.com/pkg/errors"
)

// LastBlockValidatingLoader ...
type LastBlockValidatingLoader struct {
	Loader       Loader
	Dependencies BlockDependencies
}

// LoadBlocks ...
func (loader LastBlockValidatingLoader) LoadBlocks(
	cursor interface{},
	count int,
) (
	blocks BlockGroup,
	nextCursor interface{},
	err error,
) {
	blocks, nextCursor, err = loader.Loader.LoadBlocks(cursor, count)
	if err != nil {
		return nil, nil, err
	}
	if len(blocks) == 0 {
		return blocks, nextCursor, nil
	}

	nextBlocks, _, err := loader.Loader.LoadBlocks(nextCursor, count)
	if err != nil {
		const message = "unable to preload the next blocks " +
			"corresponding to cursor %v (next cursor %v)"
		return nil, nil, errors.Wrapf(err, message, cursor, nextCursor)
	}

	var prevBlock *Block
	var validationMode ValidationMode
	if len(nextBlocks) == 0 {
		validationMode = AsFullBlockchain
	} else {
		prevBlock = &nextBlocks[0]
		validationMode = AsBlockchainChunk
	}
	err = blocks.IsLastBlockValid(prevBlock, validationMode, loader.Dependencies)
	if err != nil {
		const message = "the last block of the blocks corresponding to cursor %v " +
			"is not valid"
		return nil, nil, errors.Wrapf(err, message, cursor)
	}

	return blocks, nextCursor, nil
}
