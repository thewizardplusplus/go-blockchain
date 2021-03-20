package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBlockchain(test *testing.T) {
	type args struct {
		genesisBlockData Hasher
		dependencies     Dependencies
	}

	for _, data := range []struct {
		name          string
		args          args
		wantLastBlock Block
		wantErr       assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			gotBlockchain, gotErr :=
				NewBlockchain(data.args.genesisBlockData, data.args.dependencies)

			mock.AssertExpectationsForObjects(
				test,
				data.args.dependencies.Proofer,
				data.args.dependencies.Storage,
			)
			data.wantErr(test, gotErr)

			if gotBlockchain != nil {
				mock.AssertExpectationsForObjects(test, gotBlockchain.lastBlock.Data)
				assert.Equal(test, data.wantLastBlock, gotBlockchain.lastBlock)
			}
		})
	}
}
