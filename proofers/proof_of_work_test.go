package proofers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestProofOfWork_Hash(test *testing.T) {
	data := new(MockHasher)
	data.On("Hash").Return("hash")

	proofer := ProofOfWork{TargetBit: 248}
	hash := proofer.Hash(blockchain.Block{
		Timestamp: clock(),
		Data:      data,
		PrevHash:  "previous hash",
	})

	wantedHash :=
		"315:0093bb88b062fb387b240d14d862365f1cda9c0cda6140f19dd84deca2e91bc5"
	mock.AssertExpectationsForObjects(test, data)
	assert.Equal(test, wantedHash, hash)
}
