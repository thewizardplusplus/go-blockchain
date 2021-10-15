package loading

import (
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
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
		loader         blockchain.Loader
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
		wantBlocks         blockchain.BlockGroup
		wantNextCursor     interface{}
		wantErr            assert.ErrorAssertionFunc
	}{
		{
			name: "success with the not memoized request",
			fields: fields{
				loader: func() blockchain.Loader {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockStringer),
							Hash:      "next hash",
							PrevHash:  "hash",
						},
						{
							Timestamp: clock(),
							Data:      new(MockStringer),
							Hash:      "hash",
							PrevHash:  "previous hash",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)

					return loader
				}(),
				loadingResults: new(sync.Map),
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantLoadingResults: []memoizedRecord{
				{
					key: loadingParameters{
						cursor: "cursor-one",
						count:  23,
					},
					value: loadingResult{
						blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockStringer),
								Hash:      "next hash",
								PrevHash:  "hash",
							},
							{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "hash",
								PrevHash:  "previous hash",
							},
						},
						nextCursor: "cursor-two",
					},
				},
			},
			wantBlocks: blockchain.BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockStringer),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: clock(),
					Data:      new(MockStringer),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			wantNextCursor: "cursor-two",
			wantErr:        assert.NoError,
		},
		{
			name: "success with the memoized request",
			fields: fields{
				loader: new(MockLoader),
				loadingResults: func() *sync.Map {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockStringer),
							Hash:      "next hash",
							PrevHash:  "hash",
						},
						{
							Timestamp: clock(),
							Data:      new(MockStringer),
							Hash:      "hash",
							PrevHash:  "previous hash",
						},
					}

					parameters := loadingParameters{cursor: "cursor-one", count: 23}
					results := loadingResult{blocks: blocks, nextCursor: "cursor-two"}

					loadingResults := new(sync.Map)
					loadingResults.Store(parameters, results)

					return loadingResults
				}(),
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantLoadingResults: []memoizedRecord{
				{
					key: loadingParameters{
						cursor: "cursor-one",
						count:  23,
					},
					value: loadingResult{
						blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockStringer),
								Hash:      "next hash",
								PrevHash:  "hash",
							},
							{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "hash",
								PrevHash:  "previous hash",
							},
						},
						nextCursor: "cursor-two",
					},
				},
			},
			wantBlocks: blockchain.BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockStringer),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: clock(),
					Data:      new(MockStringer),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			wantNextCursor: "cursor-two",
			wantErr:        assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				loader: func() blockchain.Loader {
					loader := new(MockLoader)
					loader.
						On("LoadBlocks", "cursor-one", 23).
						Return(nil, nil, iotest.ErrTimeout)

					return loader
				}(),
				loadingResults: new(sync.Map),
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantLoadingResults: nil,
			wantBlocks:         nil,
			wantNextCursor:     nil,
			wantErr:            assert.Error,
		},
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
