package proofers

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/thewizardplusplus/go-blockchain"
)

// ProofOfWork ...
type ProofOfWork struct {
	TargetBit int
}

// Hash ...
func (proofer ProofOfWork) Hash(block blockchain.Block) string {
	var nonce big.Int
	var hash []byte
	target := makeTarget(proofer.TargetBit)
	for {
		data := block.MergedData() + nonce.String() + strconv.Itoa(proofer.TargetBit)
		hash = makeHash(data)
		if isHashFitTarget(hash, target) {
			break
		}

		nonce.Add(&nonce, big.NewInt(1)) // nonce += 1
	}

	return fmt.Sprintf("%d:%s:%x", proofer.TargetBit, &nonce, hash)
}

// Validate ...
func (proofer ProofOfWork) Validate(block blockchain.Block) bool {
	hashParts := strings.SplitN(block.Hash, ":", 3)
	if len(hashParts) != 3 {
		return false
	}

	targetBitAsStr := hashParts[0]
	targetBit, err := strconv.Atoi(targetBitAsStr)
	if err != nil {
		return false
	}

	target := makeTarget(targetBit)

	nonceAsStr := hashParts[1]
	data := block.MergedData() + nonceAsStr + targetBitAsStr
	hash := makeHash(data)
	return isHashFitTarget(hash, target)
}
