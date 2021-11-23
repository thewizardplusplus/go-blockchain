package loading

import (
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestLastBlockValidatingLoader_LoadBlocks(test *testing.T) {
	type fields struct {
		Loader  blockchain.Loader
		Proofer blockchain.Proofer
	}
	type args struct {
		cursor interface{}
		count  int
	}

	for _, data := range []struct {
		name           string
		fields         fields
		args           args
		wantBlocks     blockchain.BlockGroup
		wantNextCursor interface{}
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "success without blocks",
			fields: fields{
				Loader: func() blockchain.Loader {
					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(nil, "cursor-two", nil)

					return loader
				}(),
				Proofer: new(MockProofer),
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks:     nil,
			wantNextCursor: "cursor-two",
			wantErr:        assert.NoError,
		},
		{
			name: "success with blocks and without next blocks",
			fields: fields{
				Loader: func() blockchain.Loader {
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
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)
					loader.On("LoadBlocks", "cursor-two", 23).Return(nil, "cursor-three", nil)

					return loader
				}(),
				Proofer: func() blockchain.Proofer {
					proofer := new(MockProofer)
					proofer.
						On("Validate", blockchain.Block{
							Timestamp: clock(),
							Data:      new(MockStringer),
							Hash:      "hash",
							PrevHash:  "",
						}).
						Return(nil)

					return proofer
				}(),
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
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
					PrevHash:  "",
				},
			},
			wantNextCursor: "cursor-two",
			wantErr:        assert.NoError,
		},
		{
			name: "success with blocks and next blocks",
			fields: fields{
				Loader: func() blockchain.Loader {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(3 * time.Hour),
							Data:      new(MockStringer),
							Hash:      "hash #4",
							PrevHash:  "hash #3",
						},
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockStringer),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
					}
					nextBlocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockStringer),
							Hash:      "hash #2",
							PrevHash:  "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockStringer),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)
					loader.
						On("LoadBlocks", "cursor-two", 23).
						Return(nextBlocks, "cursor-three", nil)

					return loader
				}(),
				Proofer: func() blockchain.Proofer {
					proofer := new(MockProofer)
					proofer.
						On("Validate", blockchain.Block{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockStringer),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						}).
						Return(nil)

					return proofer
				}(),
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks: blockchain.BlockGroup{
				{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			wantNextCursor: "cursor-two",
			wantErr:        assert.NoError,
		},
		{
			name: "error with block loading",
			fields: fields{
				Loader: func() blockchain.Loader {
					loader := new(MockLoader)
					loader.
						On("LoadBlocks", "cursor-one", 23).
						Return(nil, "", iotest.ErrTimeout)

					return loader
				}(),
				Proofer: new(MockProofer),
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks:     nil,
			wantNextCursor: nil,
			wantErr:        assert.Error,
		},
		{
			name: "error with next block loading",
			fields: fields{
				Loader: func() blockchain.Loader {
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
					loader.
						On("LoadBlocks", "cursor-two", 23).
						Return(nil, nil, iotest.ErrTimeout)

					return loader
				}(),
				Proofer: new(MockProofer),
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks:     nil,
			wantNextCursor: nil,
			wantErr:        assert.Error,
		},
		{
			name: "error with block validating and without next blocks",
			fields: fields{
				Loader: func() blockchain.Loader {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockStringer),
							Hash:      "next hash",
							PrevHash:  "hash",
						},
						{
							Timestamp: time.Time{},
							Data:      new(MockStringer),
							Hash:      "hash",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)
					loader.On("LoadBlocks", "cursor-two", 23).Return(nil, "cursor-three", nil)

					return loader
				}(),
				Proofer: new(MockProofer),
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks:     nil,
			wantNextCursor: nil,
			wantErr:        assert.Error,
		},
		{
			name: "error with block validating and next blocks",
			fields: fields{
				Loader: func() blockchain.Loader {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(3 * time.Hour),
							Data:      new(MockStringer),
							Hash:      "hash #4",
							PrevHash:  "hash #3",
						},
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockStringer),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
					}
					nextBlocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockStringer),
							Hash:      "incorrect hash #2",
							PrevHash:  "incorrect hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockStringer),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)
					loader.
						On("LoadBlocks", "cursor-two", 23).
						Return(nextBlocks, "cursor-three", nil)

					return loader
				}(),
				Proofer: new(MockProofer),
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks:     nil,
			wantNextCursor: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			loader := LastBlockValidatingLoader{
				Loader:  data.fields.Loader,
				Proofer: data.fields.Proofer,
			}
			gotBlocks, gotNextCursor, gotErr :=
				loader.LoadBlocks(data.args.cursor, data.args.count)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.Loader,
				data.fields.Proofer,
			)
			assert.Equal(test, data.wantBlocks, gotBlocks)
			assert.Equal(test, data.wantNextCursor, gotNextCursor)
			data.wantErr(test, gotErr)
		})
	}
}
