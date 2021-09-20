package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGroupStorageWrapper_StoreBlockGroup(test *testing.T) {
	type fields struct {
		Storage Storage
	}
	type args struct {
		blocks BlockGroup
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			wrapper := GroupStorageWrapper{
				Storage: data.fields.Storage,
			}
			gotErr := wrapper.StoreBlockGroup(data.args.blocks)

			mock.AssertExpectationsForObjects(test, data.fields.Storage)
			for _, block := range data.args.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			data.wantErr(test, gotErr)
		})
	}
}
