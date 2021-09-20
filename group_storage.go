package blockchain

import (
	"github.com/pkg/errors"
)

// GroupStorage ...
type GroupStorage interface {
	Storage

	StoreBlockGroup(blocks BlockGroup) error
}

// GroupStorageWrapper ...
type GroupStorageWrapper struct {
	Storage
}

// StoreBlockGroup ...
func (wrapper GroupStorageWrapper) StoreBlockGroup(blocks BlockGroup) error {
	for index, block := range blocks {
		if err := wrapper.StoreBlock(block); err != nil {
			return errors.Wrapf(err, "unable to store block #%d", index)
		}
	}

	return nil
}
