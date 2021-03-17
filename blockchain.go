package blockchain

// Storage ...
type Storage interface {
	LoadLastBlock() (Block, error)
	StoreBlock(block Block) error
}
