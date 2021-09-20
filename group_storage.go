package blockchain

// GroupStorage ...
type GroupStorage interface {
	Storage

	StoreBlockGroup(blocks BlockGroup) error
}

// GroupStorageWrapper ...
type GroupStorageWrapper struct {
	Storage
}
