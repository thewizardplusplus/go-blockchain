package loading

import (
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestChunkValidatingLoader_LoadBlocks(test *testing.T) {
	type fields struct {
		Loader       blockchain.Loader
		Dependencies blockchain.BlockDependencies
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
				Dependencies: blockchain.BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
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
			name: "success with blocks",
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

					return loader
				}(),
				Dependencies: blockchain.BlockDependencies{
					Clock: clock,
					Proofer: func() blockchain.Proofer {
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

						proofer := new(MockProofer)
						for _, block := range blocks {
							proofer.On("Validate", block).Return(nil)
						}

						return proofer
					}(),
				},
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
					PrevHash:  "previous hash",
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
				Dependencies: blockchain.BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
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
			name: "error with block validating",
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

					return loader
				}(),
				Dependencies: blockchain.BlockDependencies{
					Clock: clock,
					Proofer: func() blockchain.Proofer {
						block := blockchain.Block{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockStringer),
							Hash:      "next hash",
							PrevHash:  "hash",
						}

						proofer := new(MockProofer)
						proofer.On("Validate", block).Return(iotest.ErrTimeout)

						return proofer
					}(),
				},
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
