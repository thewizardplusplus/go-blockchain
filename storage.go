package blockchain

//go:generate mockery --name=Storage --inpackage --case=underscore --testonly

// Storage ...
type Storage interface {
	LoadLastBlock() (Block, error)
	StoreBlock(block Block) error
}
