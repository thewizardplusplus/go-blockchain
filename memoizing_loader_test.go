package blockchain

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewMemoizingLoader(test *testing.T) {
	loader := new(MockLoader)
	memoizingLoader := NewMemoizingLoader(loader)

	mock.AssertExpectationsForObjects(test, loader)
	assert.Equal(test, loader, memoizingLoader.loader)
	assert.NotNil(test, memoizingLoader.loadingResults)
}

func TestMemoizingLoader_LoadBlocks(test *testing.T) {
	type fields struct {
		loader         Loader
		loadingResults *sync.Map
	}
	type args struct {
		cursor interface{}
		count  int
	}
	type memoizedRecord struct {
		key   loadingParameters
		value loadingResult
	}

	for _, data := range []struct {
		name               string
		fields             fields
		args               args
		wantLoadingResults []memoizedRecord
		wantBlocks         BlockGroup
		wantNextCursor     interface{}
		wantErr            assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			loader := MemoizingLoader{
				loader:         data.fields.loader,
				loadingResults: data.fields.loadingResults,
			}
			gotBlocks, gotNextCursor, gotErr :=
				loader.LoadBlocks(data.args.cursor, data.args.count)

			var gotLoadingResults []memoizedRecord
			data.fields.loadingResults.
				Range(func(key interface{}, value interface{}) bool {
					gotLoadingResults = append(gotLoadingResults, memoizedRecord{
						key:   key.(loadingParameters),
						value: value.(loadingResult),
					})

					return true
				})

			mock.AssertExpectationsForObjects(test, data.fields.loader)
			assert.Equal(test, data.wantLoadingResults, gotLoadingResults)
			assert.Equal(test, data.wantBlocks, gotBlocks)
			assert.Equal(test, data.wantNextCursor, gotNextCursor)
			data.wantErr(test, gotErr)
		})
	}
}
