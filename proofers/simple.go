package proofers

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/thewizardplusplus/go-blockchain"
)

// Simple ...
type Simple struct{}

// Hash ...
func (proofer Simple) Hash(block blockchain.Block) string {
	data := block.Timestamp.String() + block.Data.Hash() + block.PrevHash
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
