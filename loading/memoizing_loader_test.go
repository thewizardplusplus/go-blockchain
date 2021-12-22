package loading

import (
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestNewMemoizingLoader(test *testing.T) {
	maximalCacheSize := int(1e6)
	loader := new(MockLoader)
	memoizingLoader := NewMemoizingLoader(maximalCacheSize, loader)

	mock.AssertExpectationsForObjects(test, loader)
	assert.Equal(test, loader, memoizingLoader.loader)
	assert.Equal(
		test,
		NewLRUCache(maximalCacheSize),
		memoizingLoader.loadingResults,
	)
}

func TestMemoizingLoader_LoadBlocks(test *testing.T) {
	type fields struct {
		loader         blockchain.Loader
		loadingResults LRUCache
	}
	type args struct {
		cursor interface{}
		count  int
	}

	for _, data := range []struct {
		name               string
		fields             fields
		args               args
		wantLoadingResults LRUCache
		wantBlocks         blockchain.BlockGroup
		wantNextCursor     interface{}
		wantErr            assert.ErrorAssertionFunc
	}{
		{
			name: "success with the memoized request",
			fields: fields{
				loader: new(MockLoader),
				loadingResults: func() LRUCache {
					loadingResults := NewLRUCache(10)
					loadingResults.Set(
						Parameters{Cursor: "cursor #1", Count: 2},
						Results{
							Blocks: blockchain.BlockGroup{
								{
									Timestamp: clock(),
									Data:      new(MockData),
									Hash:      "hash #1",
									PrevHash:  "",
								},
								{
									Timestamp: clock().Add(time.Hour),
									Data:      new(MockData),
									Hash:      "hash #2",
									PrevHash:  "hash #1",
								},
							},
							NextCursor: "cursor #2",
						},
					)
					loadingResults.Set(
						Parameters{Cursor: "cursor #2", Count: 2},
						Results{
							Blocks: blockchain.BlockGroup{
								{
									Timestamp: clock().Add(2 * time.Hour),
									Data:      new(MockData),
									Hash:      "hash #3",
									PrevHash:  "hash #2",
								},
								{
									Timestamp: clock().Add(3 * time.Hour),
									Data:      new(MockData),
									Hash:      "hash #4",
									PrevHash:  "hash #3",
								},
							},
							NextCursor: "cursor #3",
						},
					)

					return loadingResults
				}(),
			},
			args: args{
				cursor: "cursor #2",
				count:  2,
			},
			wantLoadingResults: func() LRUCache {
				loadingResults := NewLRUCache(10)
				loadingResults.Set(
					Parameters{Cursor: "cursor #1", Count: 2},
					Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				)
				loadingResults.Set(
					Parameters{Cursor: "cursor #2", Count: 2},
					Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(3 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #4",
								PrevHash:  "hash #3",
							},
						},
						NextCursor: "cursor #3",
					},
				)

				return loadingResults
			}(),
			wantBlocks: blockchain.BlockGroup{
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
				{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
			},
			wantNextCursor: "cursor #3",
			wantErr:        assert.NoError,
		},
		{
			name: "success with the not memoized request",
			fields: fields{
				loader: func() blockchain.Loader {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(4 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #5",
							PrevHash:  "hash #4",
						},
						{
							Timestamp: clock().Add(5 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #6",
							PrevHash:  "hash #5",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor #3", 2).Return(blocks, "cursor #4", nil)

					return loader
				}(),
				loadingResults: func() LRUCache {
					loadingResults := NewLRUCache(10)
					loadingResults.Set(
						Parameters{Cursor: "cursor #1", Count: 2},
						Results{
							Blocks: blockchain.BlockGroup{
								{
									Timestamp: clock(),
									Data:      new(MockData),
									Hash:      "hash #1",
									PrevHash:  "",
								},
								{
									Timestamp: clock().Add(time.Hour),
									Data:      new(MockData),
									Hash:      "hash #2",
									PrevHash:  "hash #1",
								},
							},
							NextCursor: "cursor #2",
						},
					)
					loadingResults.Set(
						Parameters{Cursor: "cursor #2", Count: 2},
						Results{
							Blocks: blockchain.BlockGroup{
								{
									Timestamp: clock().Add(2 * time.Hour),
									Data:      new(MockData),
									Hash:      "hash #3",
									PrevHash:  "hash #2",
								},
								{
									Timestamp: clock().Add(3 * time.Hour),
									Data:      new(MockData),
									Hash:      "hash #4",
									PrevHash:  "hash #3",
								},
							},
							NextCursor: "cursor #3",
						},
					)

					return loadingResults
				}(),
			},
			args: args{
				cursor: "cursor #3",
				count:  2,
			},
			wantLoadingResults: func() LRUCache {
				loadingResults := NewLRUCache(10)
				loadingResults.Set(
					Parameters{Cursor: "cursor #1", Count: 2},
					Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				)
				loadingResults.Set(
					Parameters{Cursor: "cursor #2", Count: 2},
					Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(3 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #4",
								PrevHash:  "hash #3",
							},
						},
						NextCursor: "cursor #3",
					},
				)
				loadingResults.Set(
					Parameters{Cursor: "cursor #3", Count: 2},
					Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(4 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #5",
								PrevHash:  "hash #4",
							},
							{
								Timestamp: clock().Add(5 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #6",
								PrevHash:  "hash #5",
							},
						},
						NextCursor: "cursor #4",
					},
				)

				return loadingResults
			}(),
			wantBlocks: blockchain.BlockGroup{
				{
					Timestamp: clock().Add(4 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #5",
					PrevHash:  "hash #4",
				},
				{
					Timestamp: clock().Add(5 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #6",
					PrevHash:  "hash #5",
				},
			},
			wantNextCursor: "cursor #4",
			wantErr:        assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				loader: func() blockchain.Loader {
					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor #1", 2).Return(nil, nil, iotest.ErrTimeout)

					return loader
				}(),
				loadingResults: NewLRUCache(10),
			},
			args: args{
				cursor: "cursor #1",
				count:  2,
			},
			wantLoadingResults: NewLRUCache(10),
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

			mock.AssertExpectationsForObjects(test, data.fields.loader)
			assert.Equal(test, data.wantLoadingResults, loader.loadingResults)
			assert.Equal(test, data.wantBlocks, gotBlocks)
			assert.Equal(test, data.wantNextCursor, gotNextCursor)
			data.wantErr(test, gotErr)
		})
	}
}
