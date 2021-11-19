package blockchain

import (
	"fmt"
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBlockchain(test *testing.T) {
	type args struct {
		genesisBlockData fmt.Stringer
		dependencies     Dependencies
	}

	for _, data := range []struct {
		name          string
		args          args
		wantLastBlock Block
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "success with a nonempty storage",
			args: args{
				genesisBlockData: new(MockStringer),
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock:   clock,
						Proofer: new(MockProofer),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.
							On("LoadLastBlock").
							Return(
								Block{
									Timestamp: clock(),
									Data:      new(MockStringer),
									Hash:      "hash",
									PrevHash:  "previous hash",
								},
								nil,
							)

						return storage
					}(),
				},
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with an empty storage",
			args: args{
				genesisBlockData: new(MockStringer),
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock: clock,
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.
								On("Hash", Block{
									Timestamp: clock(),
									Data:      new(MockStringer),
									PrevHash:  "",
								}).
								Return("hash")

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.On("LoadLastBlock").Return(Block{}, ErrEmptyStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "hash",
								PrevHash:  "",
							}).
							Return(nil)

						return storage
					}(),
				},
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error on last block loading",
			args: args{
				genesisBlockData: new(MockStringer),
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock:   clock,
						Proofer: new(MockProofer),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.On("LoadLastBlock").Return(Block{}, iotest.ErrTimeout)

						return storage
					}(),
				},
			},
			wantLastBlock: Block{},
			wantErr:       assert.Error,
		},
		{
			name: "error with an empty storage",
			args: args{
				genesisBlockData: nil,
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock:   clock,
						Proofer: new(MockProofer),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.On("LoadLastBlock").Return(Block{}, ErrEmptyStorage)

						return storage
					}(),
				},
			},
			wantLastBlock: Block{},
			wantErr:       assert.Error,
		},
		{
			name: "error on genesis block storing",
			args: args{
				genesisBlockData: new(MockStringer),
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock: clock,
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.
								On("Hash", Block{
									Timestamp: clock(),
									Data:      new(MockStringer),
									PrevHash:  "",
								}).
								Return("hash")

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.On("LoadLastBlock").Return(Block{}, ErrEmptyStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "hash",
								PrevHash:  "",
							}).
							Return(iotest.ErrTimeout)

						return storage
					}(),
				},
			},
			wantLastBlock: Block{},
			wantErr:       assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotBlockchain, gotErr :=
				NewBlockchain(data.args.genesisBlockData, data.args.dependencies)

			if data.args.genesisBlockData != nil {
				mock.AssertExpectationsForObjects(test, data.args.genesisBlockData)
			}
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

func TestBlockchain_LoadBlocks(test *testing.T) {
	type fields struct {
		dependencies Dependencies
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
		{
			name: "success",
			fields: fields{
				dependencies: Dependencies{
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.
							On("LoadBlocks", 1, 2).
							Return(
								BlockGroup{
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
								3,
								nil,
							)

						return storage
					}(),
				},
			},
			args: args{
				cursor: 1,
				count:  2,
			},
			wantBlocks: BlockGroup{
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
			wantNextCursor: 3,
			wantErr:        assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				dependencies: Dependencies{
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.On("LoadBlocks", 1, 2).Return(nil, nil, iotest.ErrTimeout)

						return storage
					}(),
				},
			},
			args: args{
				cursor: 1,
				count:  2,
			},
			wantBlocks:     nil,
			wantNextCursor: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			blockchain := Blockchain{
				dependencies: data.fields.dependencies,
			}
			gotBlocks, gotNextCursor, gotErr :=
				blockchain.LoadBlocks(data.args.cursor, data.args.count)

			assert.Equal(test, data.wantBlocks, gotBlocks)
			assert.Equal(test, data.wantNextCursor, gotNextCursor)
			data.wantErr(test, gotErr)
		})
	}
}

func TestBlockchain_AddBlock(test *testing.T) {
	type fields struct {
		dependencies Dependencies
		lastBlock    Block
	}
	type args struct {
		data fmt.Stringer
	}

	for _, data := range []struct {
		name          string
		fields        fields
		args          args
		wantLastBlock Block
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock: clock,
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.
								On("Hash", Block{
									Timestamp: clock(),
									Data:      new(MockStringer),
									PrevHash:  "hash",
								}).
								Return("next hash")

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "next hash",
								PrevHash:  "hash",
							}).
							Return(nil)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockStringer),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				data: new(MockStringer),
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "next hash",
				PrevHash:  "hash",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock: clock,
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.
								On("Hash", Block{
									Timestamp: clock(),
									Data:      new(MockStringer),
									PrevHash:  "hash",
								}).
								Return("next hash")

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "next hash",
								PrevHash:  "hash",
							}).
							Return(iotest.ErrTimeout)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockStringer),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				data: new(MockStringer),
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(t *testing.T) {
			blockchain := &Blockchain{
				dependencies: data.fields.dependencies,
				lastBlock:    data.fields.lastBlock,
			}
			gotErr := blockchain.AddBlock(data.args.data)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.dependencies.Proofer,
				data.fields.dependencies.Storage,
				data.fields.lastBlock.Data,
				data.args.data,
				blockchain.lastBlock.Data,
			)
			assert.Equal(test, data.wantLastBlock, blockchain.lastBlock)
			data.wantErr(test, gotErr)
		})
	}
}
