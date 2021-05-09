package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBlockGroup_IsValid(test *testing.T) {
	type args struct {
		prependedChunk BlockGroup
		validationMode ValidationMode
		dependencies   BlockDependencies
	}

	for _, data := range []struct {
		name   string
		blocks BlockGroup
		args   args
		want   assert.BoolAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.blocks.IsValid(
				data.args.prependedChunk,
				data.args.validationMode,
				data.args.dependencies,
			)

			for _, block := range data.args.prependedChunk {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			for _, block := range data.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			mock.AssertExpectationsForObjects(test, data.args.dependencies.Proofer)
			data.want(test, got)
		})
	}
}
