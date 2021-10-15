package proofers

import (
	"encoding/hex"

	"github.com/pkg/errors"
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
func (proofer Simple) Validate(block blockchain.Block) error {
	if block.Hash != proofer.Hash(block) {
		return errors.New("the hash is not valid")
	}

	return nil
}
