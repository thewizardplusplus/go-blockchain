package blockchain

// ChunkValidatingLoader ...
type ChunkValidatingLoader struct {
	Loader       Loader
	Dependencies BlockDependencies
}
