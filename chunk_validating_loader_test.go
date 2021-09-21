package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChunkValidatingLoader_LoadBlocks(test *testing.T) {
	type fields struct {
		Loader       Loader
		Dependencies BlockDependencies
	}
	type args struct {
		cursor interface{}
		count  int
	}

	for _, data := range []struct {
		name           string
		fields         fields
		args           args
		wantBlocks     BlockGroup
		wantNextCursor interface{}
		wantErr        assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			loader := ChunkValidatingLoader{
				Loader:       data.fields.Loader,
				Dependencies: data.fields.Dependencies,
			}
			gotBlocks, gotNextCursor, gotErr :=
				loader.LoadBlocks(data.args.cursor, data.args.count)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.Loader,
				data.fields.Dependencies.Proofer,
			)
			assert.Equal(test, data.wantBlocks, gotBlocks)
			assert.Equal(test, data.wantNextCursor, gotNextCursor)
			data.wantErr(test, gotErr)
		})
	}
}
