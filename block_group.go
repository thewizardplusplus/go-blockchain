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
	validationMode ValidationMode,
	dependencies BlockDependencies,
) bool {
	lastIndex := len(blocks) - 1
	for index, block := range blocks[:lastIndex] {
		prevBlock := &blocks[index+1]
		if !block.IsValid(prevBlock, dependencies) {
			return false
		}
	}

	var lastBlockValidator func(block Block) bool
	switch validationMode {
	case AsFullBlockchain:
		lastBlockValidator =
			func(block Block) bool { return block.IsValidGenesisBlock(dependencies) }
	case AsBlockchainChunk:
		lastBlockValidator =
			func(block Block) bool { return block.IsValid(nil, dependencies) }
	}
	if !lastBlockValidator(blocks[lastIndex]) {
		return false
	}

	return true
}
