package storing

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-blockchain"
)

// GroupStorageWrapper ...
type GroupStorageWrapper struct {
	blockchain.Storage
}

// StoreBlockGroup ...
func (wrapper GroupStorageWrapper) StoreBlockGroup(
	blocks blockchain.BlockGroup,
) error {
	for index, block := range blocks {
		if err := wrapper.StoreBlock(block); err != nil {
			return errors.Wrapf(err, "unable to store block #%d", index)
		}
	}

	return nil
}
