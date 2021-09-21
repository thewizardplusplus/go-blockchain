package blockchain

// ValidationMode ...
type ValidationMode int

// ...
const (
	AsFullBlockchain ValidationMode = iota
	AsBlockchainChunk
)

// BlockGroup ...
type BlockGroup []Block

// IsValid ...
func (blocks BlockGroup) IsValid(
	prependedChunk BlockGroup,
	validationMode ValidationMode,
	dependencies BlockDependencies,
) bool {
	if len(blocks) == 0 {
		return true
	}

	if len(prependedChunk) != 0 {
		prevBlock, validationMode := &blocks[0], AsBlockchainChunk
		if !prependedChunk.IsLastBlockValid(prevBlock, validationMode, dependencies) {
			return false
		}
	}

	for index, block := range blocks[:len(blocks)-1] {
		prevBlock := &blocks[index+1]
		if !block.IsValid(prevBlock, dependencies) {
			return false
		}
	}

	return blocks.IsLastBlockValid(nil, validationMode, dependencies)
}

// IsLastBlockValid ...
func (blocks BlockGroup) IsLastBlockValid(
	prevBlock *Block,
	validationMode ValidationMode,
	dependencies BlockDependencies,
) bool {
	var isValid bool
	lastBlock := blocks[len(blocks)-1]
	switch validationMode {
	case AsFullBlockchain:
		isValid = lastBlock.IsValidGenesisBlock(dependencies)
	case AsBlockchainChunk:
		isValid = lastBlock.IsValid(prevBlock, dependencies)
	}

	return isValid
}
