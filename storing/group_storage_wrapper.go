package storing

import (
	"fmt"

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
			return fmt.Errorf("unable to store block #%d: %w", index, err)
		}
	}

	return nil
}

// DeleteBlockGroup ...
func (wrapper GroupStorageWrapper) DeleteBlockGroup(
	blocks blockchain.BlockGroup,
) error {
	for index, block := range blocks {
		if err := wrapper.DeleteBlock(block); err != nil {
			return fmt.Errorf("unable to delete block #%d: %w", index, err)
		}
	}

	return nil
}
