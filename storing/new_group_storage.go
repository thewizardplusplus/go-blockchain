package storing

import (
	"github.com/thewizardplusplus/go-blockchain"
)

// NewGroupStorage ...
func NewGroupStorage(storage blockchain.Storage) blockchain.GroupStorage {
	groupStorage, ok := storage.(blockchain.GroupStorage)
	if ok {
		return groupStorage
	}

	return GroupStorageWrapper{Storage: storage}
}
