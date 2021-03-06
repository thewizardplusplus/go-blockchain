package proofers

import (
	"crypto/sha256"
	"math/big"
)

func makeHash(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

func makeTarget(targetBit int) *big.Int {
	target := big.NewInt(0)
	target.SetBit(target, targetBit, 1)

	return target
}

func isHashFitTarget(hash []byte, target *big.Int) bool {
	hashAsInt := big.NewInt(0)
	hashAsInt.SetBytes(hash)

	return hashAsInt.Cmp(target) == -1 // is less
}
