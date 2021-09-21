package blockchain

// LastBlockValidatingLoader ...
type LastBlockValidatingLoader struct {
	Loader       Loader
	Dependencies BlockDependencies
}
