package proofers

import (
	"encoding/hex"

	"github.com/thewizardplusplus/go-blockchain"
)

// Simple ...
type Simple struct{}

// Hash ...
func (proofer Simple) Hash(block blockchain.Block) string {
	data := block.MergedData()
	hash := makeHash(data)
	return hex.EncodeToString(hash)
}

// Validate ...
func (proofer Simple) Validate(block blockchain.Block) bool {
	return block.Hash == proofer.Hash(block)
}
