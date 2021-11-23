package loading

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-blockchain"
)

// LastBlockValidatingLoader ...
type LastBlockValidatingLoader struct {
	Loader  blockchain.Loader
	Proofer blockchain.Proofer
}

// LoadBlocks ...
func (loader LastBlockValidatingLoader) LoadBlocks(
	cursor interface{},
	count int,
) (
	blocks blockchain.BlockGroup,
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

	var prevBlock *blockchain.Block
	var validationMode blockchain.ValidationMode
	if len(nextBlocks) == 0 {
		validationMode = blockchain.AsFullBlockchain
	} else {
		prevBlock = &nextBlocks[0]
		validationMode = blockchain.AsBlockchainChunk
	}
	err = blocks.IsLastBlockValid(prevBlock, validationMode, loader.Proofer)
	if err != nil {
		const message = "the last block of the blocks corresponding to cursor %v " +
			"is not valid"
		return nil, nil, errors.Wrapf(err, message, cursor)
	}

	return blocks, nextCursor, nil
}
