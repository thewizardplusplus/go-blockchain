package proofers

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

	"github.com/thewizardplusplus/go-blockchain"
)

// ProofOfWork ...
type ProofOfWork struct {
	TargetBit int
}

// Hash ...
func (proofer ProofOfWork) Hash(block blockchain.Block) string {
	var target big.Int
	target.SetBit(&target, proofer.TargetBit, 1)

	var nonce big.Int
	var hash [sha256.Size]byte
	for {
		data := block.MergedData() + nonce.String()
		hash = sha256.Sum256([]byte(data))

		hashAsInt := big.NewInt(0)
		hashAsInt.SetBytes(hash[:])

		if hashAsInt.Cmp(&target) == -1 /* is less */ {
			break
		}

		nonce.Add(&nonce, big.NewInt(1)) // nonce += 1
	}

	return fmt.Sprintf("%s:%x", &nonce, hash)
}

// Validate ...
func (proofer ProofOfWork) Validate(block blockchain.Block) bool {
	var target big.Int
	target.SetBit(&target, proofer.TargetBit, 1)

	hashParts := strings.SplitN(block.Hash, ":", 2)
	if len(hashParts) != 2 {
		return false
	}

	nonceAsStr := hashParts[0]
	data := block.MergedData() + nonceAsStr
	hash := sha256.Sum256([]byte(data))

	hashAsInt := big.NewInt(0)
	hashAsInt.SetBytes(hash[:])

	return hashAsInt.Cmp(&target) == -1 /* is less */
}
