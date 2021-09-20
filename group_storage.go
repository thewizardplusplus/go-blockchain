package blockchain

import (
	"github.com/pkg/errors"
)

//go:generate mockery --name=GroupStorage --inpackage --case=underscore --testonly

// GroupStorage ...
type GroupStorage interface {
	Storage

	StoreBlockGroup(blocks BlockGroup) error
}

// NewGroupStorage ...
func NewGroupStorage(storage Storage) GroupStorage {
	groupStorage, ok := storage.(GroupStorage)
	if ok {
		return groupStorage
	}

	return GroupStorageWrapper{Storage: storage}
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
