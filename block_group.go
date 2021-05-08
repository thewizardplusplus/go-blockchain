package blockchain

// BlockGroup ...
type BlockGroup []Block

// IsValid ...
func (blocks BlockGroup) IsValid(dependencies BlockDependencies) bool {
	lastIndex := len(blocks) - 1
	for index, block := range blocks[:lastIndex] {
		prevBlock := blocks[index+1]
		if !block.IsValid(prevBlock, dependencies) {
			return false
		}
	}

	if !blocks[lastIndex].IsValidGenesisBlock(dependencies) {
		return false
	}

	return true
}
