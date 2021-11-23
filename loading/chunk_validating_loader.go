package loading

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-blockchain"
)

// ChunkValidatingLoader ...
type ChunkValidatingLoader struct {
	Loader  blockchain.Loader
	Proofer blockchain.Proofer
}

// LoadBlocks ...
func (loader ChunkValidatingLoader) LoadBlocks(cursor interface{}, count int) (
	blocks blockchain.BlockGroup,
	nextCursor interface{},
	err error,
) {
	blocks, nextCursor, err = loader.Loader.LoadBlocks(cursor, count)
	if err != nil {
		return nil, nil, err
	}

	err = blocks.IsValid(nil, blockchain.AsBlockchainChunk, loader.Proofer)
	if err != nil {
		const message = "the blocks corresponding to cursor %v are not valid"
		return nil, nil, errors.Wrapf(err, message, cursor)
	}

	return blocks, nextCursor, nil
}
