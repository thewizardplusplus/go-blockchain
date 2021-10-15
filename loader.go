package blockchain

//go:generate mockery --name=Loader --inpackage --case=underscore --testonly

// Loader ...
type Loader interface {
	LoadBlocks(cursor interface{}, count int) (
		blocks BlockGroup,
		nextCursor interface{},
		err error,
	)
}
