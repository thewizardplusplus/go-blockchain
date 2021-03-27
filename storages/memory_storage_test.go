package storages

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestMemoryStorage_LoadLastBlock(test *testing.T) {
	type fields struct {
		blocks []blockchain.Block
	}

	for _, data := range []struct {
		name          string
		fields        fields
		wantLastBlock blockchain.Block
		wantErr       assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			storage := MemoryStorage{
				blocks: data.fields.blocks,
			}
			gotLastBlock, gotErr := storage.LoadLastBlock()

			for _, block := range data.fields.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			assert.Equal(test, data.wantLastBlock, gotLastBlock)
			data.wantErr(test, gotErr)
		})
	}
}
