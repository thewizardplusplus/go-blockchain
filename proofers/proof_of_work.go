package proofers

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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
func (proofer ProofOfWork) Validate(block blockchain.Block) error {
	hashParts := strings.SplitN(block.Hash, ":", 3)
	if len(hashParts) != 3 {
		return errors.New("the hash contains the invalid quantity of the parts")
	}

	targetBitAsStr := hashParts[0]
	targetBit, err := strconv.Atoi(targetBitAsStr)
	if err != nil {
		return errors.Wrap(err, "unable to parse the target bit")
	}

	target := makeTarget(targetBit)

	nonceAsStr := hashParts[1]
	data := block.MergedData() + nonceAsStr + targetBitAsStr
	hash := makeHash(data)
	if !isHashFitTarget(hash, target) {
		return errors.New("the hash does not fit the target")
	}

	return nil
}

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
