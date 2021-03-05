package proofers

import (
	"crypto/sha256"
	"fmt"
	"math/big"

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
		data := block.Timestamp.String() +
			block.Data.Hash() +
			block.PrevHash +
			nonce.String()
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
