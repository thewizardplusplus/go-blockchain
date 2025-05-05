package blockchain

import (
	"context"
	"testing"
	"testing/iotest"
	"time"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBlockchain(test *testing.T) {
	type args struct {
		genesisBlockData Data
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
				genesisBlockData: new(MockData),
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
									Data:      new(MockData),
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
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with an empty storage",
			args: args{
				genesisBlockData: new(MockData),
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock: clock,
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.
								On(
									"HashEx",
									context.Background(),
									Block{
										Timestamp: clock(),
										Data:      new(MockData),
										PrevHash:  "",
									},
								).
								Return("hash", nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.On("LoadLastBlock").Return(Block{}, ErrEmptyStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockData),
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
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error on last block loading",
			args: args{
				genesisBlockData: new(MockData),
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
				genesisBlockData: new(MockData),
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock: clock,
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.
								On(
									"HashEx",
									context.Background(),
									Block{
										Timestamp: clock(),
										Data:      new(MockData),
										PrevHash:  "",
									},
								).
								Return("hash", nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.On("LoadLastBlock").Return(Block{}, ErrEmptyStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockData),
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

func TestNewBlockchainEx(test *testing.T) {
	type args struct {
		ctx    context.Context
		params NewBlockchainExParams
	}

	for _, data := range []struct {
		name          string
		args          args
		want          assert.ValueAssertionFunc
		wantLastBlock Block
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "success/with a nonempty storage",
			args: args{
				ctx: context.Background(),
				params: NewBlockchainExParams{
					Dependencies: Dependencies{
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
										Data:      new(MockData),
										Hash:      "hash",
										PrevHash:  "previous hash",
									},
									nil,
								)

							return storage
						}(),
					},
					GenesisBlockData: mo.Some[Data](new(MockData)),
				},
			},
			want: assert.NotNil,
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/with an empty storage",
			args: args{
				ctx: context.Background(),
				params: NewBlockchainExParams{
					Dependencies: Dependencies{
						BlockDependencies: BlockDependencies{
							Clock: clock,
							Proofer: func() Proofer {
								proofer := new(MockProofer)
								proofer.
									On(
										"HashEx",
										context.Background(),
										Block{
											Timestamp: clock(),
											Data:      new(MockData),
											PrevHash:  "",
										},
									).
									Return("hash", nil)

								return proofer
							}(),
						},

						Storage: func() GroupStorage {
							storage := new(MockGroupStorage)
							storage.
								On("LoadLastBlock").
								Return(Block{}, ErrEmptyStorage)
							storage.
								On("StoreBlock", Block{
									Timestamp: clock(),
									Data:      new(MockData),
									Hash:      "hash",
									PrevHash:  "",
								}).
								Return(nil)

							return storage
						}(),
					},
					GenesisBlockData: mo.Some[Data](new(MockData)),
				},
			},
			want: assert.NotNil,
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to load the last block/regular error",
			args: args{
				ctx: context.Background(),
				params: NewBlockchainExParams{
					Dependencies: Dependencies{
						BlockDependencies: BlockDependencies{
							Clock:   clock,
							Proofer: new(MockProofer),
						},

						Storage: func() GroupStorage {
							storage := new(MockGroupStorage)
							storage.
								On("LoadLastBlock").
								Return(Block{}, iotest.ErrTimeout)

							return storage
						}(),
					},
					GenesisBlockData: mo.Some[Data](new(MockData)),
				},
			},
			want:          assert.Nil,
			wantLastBlock: Block{},
			wantErr:       assert.Error,
		},
		{
			name: "error/unable to load the last block/empty storage",
			args: args{
				ctx: context.Background(),
				params: NewBlockchainExParams{
					Dependencies: Dependencies{
						BlockDependencies: BlockDependencies{
							Clock:   clock,
							Proofer: new(MockProofer),
						},

						Storage: func() GroupStorage {
							storage := new(MockGroupStorage)
							storage.
								On("LoadLastBlock").
								Return(Block{}, ErrEmptyStorage)

							return storage
						}(),
					},
					GenesisBlockData: mo.None[Data](),
				},
			},
			want:          assert.Nil,
			wantLastBlock: Block{},
			wantErr:       assert.Error,
		},
		{
			name: "error/unable to create a new genesis block",
			args: args{
				ctx: context.Background(),
				params: NewBlockchainExParams{
					Dependencies: Dependencies{
						BlockDependencies: BlockDependencies{
							Clock: clock,
							Proofer: func() Proofer {
								proofer := new(MockProofer)
								proofer.
									On(
										"HashEx",
										context.Background(),
										Block{
											Timestamp: clock(),
											Data:      new(MockData),
											PrevHash:  "",
										},
									).
									Return("", iotest.ErrTimeout)

								return proofer
							}(),
						},

						Storage: func() GroupStorage {
							storage := new(MockGroupStorage)
							storage.
								On("LoadLastBlock").
								Return(Block{}, ErrEmptyStorage)

							return storage
						}(),
					},
					GenesisBlockData: mo.Some[Data](new(MockData)),
				},
			},
			want:          assert.Nil,
			wantLastBlock: Block{},
			wantErr:       assert.Error,
		},
		{
			name: "error/unable to store the genesis block",
			args: args{
				ctx: context.Background(),
				params: NewBlockchainExParams{
					Dependencies: Dependencies{
						BlockDependencies: BlockDependencies{
							Clock: clock,
							Proofer: func() Proofer {
								proofer := new(MockProofer)
								proofer.
									On(
										"HashEx",
										context.Background(),
										Block{
											Timestamp: clock(),
											Data:      new(MockData),
											PrevHash:  "",
										},
									).
									Return("hash", nil)

								return proofer
							}(),
						},

						Storage: func() GroupStorage {
							storage := new(MockGroupStorage)
							storage.
								On("LoadLastBlock").
								Return(Block{}, ErrEmptyStorage)
							storage.
								On("StoreBlock", Block{
									Timestamp: clock(),
									Data:      new(MockData),
									Hash:      "hash",
									PrevHash:  "",
								}).
								Return(iotest.ErrTimeout)

							return storage
						}(),
					},
					GenesisBlockData: mo.Some[Data](new(MockData)),
				},
			},
			want:          assert.Nil,
			wantLastBlock: Block{},
			wantErr:       assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewBlockchainEx(data.args.ctx, data.args.params)

			data.want(test, got)
			if got != nil {
				assert.Equal(test, data.wantLastBlock, got.lastBlock)
			}
			data.wantErr(test, err)

			if genesisBlockData, isPresent :=
				data.args.params.GenesisBlockData.Get(); isPresent {
				mock.AssertExpectationsForObjects(test, genesisBlockData)
			}
			mock.AssertExpectationsForObjects(
				test,
				data.args.params.Dependencies.Proofer,
				data.args.params.Dependencies.Storage,
			)
			if got != nil {
				mock.AssertExpectationsForObjects(test, got.lastBlock.Data)
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
										Data:      new(MockData),
										Hash:      "next hash",
										PrevHash:  "hash",
									},
									{
										Timestamp: clock(),
										Data:      new(MockData),
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
					Data:      new(MockData),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: clock(),
					Data:      new(MockData),
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
		data Data
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
								On(
									"HashEx",
									context.Background(),
									Block{
										Timestamp: clock(),
										Data:      new(MockData),
										PrevHash:  "hash",
									},
								).
								Return("next hash", nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "next hash",
								PrevHash:  "hash",
							}).
							Return(nil)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				data: new(MockData),
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
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
								On(
									"HashEx",
									context.Background(),
									Block{
										Timestamp: clock(),
										Data:      new(MockData),
										PrevHash:  "hash",
									},
								).
								Return("next hash", nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "next hash",
								PrevHash:  "hash",
							}).
							Return(iotest.ErrTimeout)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				data: new(MockData),
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
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

func TestBlockchain_AddBlockEx(test *testing.T) {
	type fields struct {
		dependencies Dependencies
		lastBlock    Block
	}
	type args struct {
		ctx  context.Context
		data Data
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
								On(
									"HashEx",
									context.Background(),
									Block{
										Timestamp: clock(),
										Data:      new(MockData),
										PrevHash:  "hash",
									},
								).
								Return("next hash", nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "next hash",
								PrevHash:  "hash",
							}).
							Return(nil)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				ctx:  context.Background(),
				data: new(MockData),
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "next hash",
				PrevHash:  "hash",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to create a new block",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock: clock,
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.
								On(
									"HashEx",
									context.Background(),
									Block{
										Timestamp: clock(),
										Data:      new(MockData),
										PrevHash:  "hash",
									},
								).
								Return("", iotest.ErrTimeout)

							return proofer
						}(),
					},
					Storage: new(MockGroupStorage),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				ctx:  context.Background(),
				data: new(MockData),
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to store the block",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock: clock,
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.
								On(
									"HashEx",
									context.Background(),
									Block{
										Timestamp: clock(),
										Data:      new(MockData),
										PrevHash:  "hash",
									},
								).
								Return("next hash", nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "next hash",
								PrevHash:  "hash",
							}).
							Return(iotest.ErrTimeout)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				ctx:  context.Background(),
				data: new(MockData),
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			blockchain := &Blockchain{
				dependencies: data.fields.dependencies,
				lastBlock:    data.fields.lastBlock,
			}
			err := blockchain.AddBlockEx(data.args.ctx, data.args.data)

			assert.Equal(test, data.wantLastBlock, blockchain.lastBlock)
			data.wantErr(test, err)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.dependencies.Proofer,
				data.fields.dependencies.Storage,
				data.fields.lastBlock.Data,
				data.args.data,
				blockchain.lastBlock.Data,
			)
		})
	}
}

func TestBlockchain_Merge(test *testing.T) {
	type fields struct {
		dependencies Dependencies
		lastBlock    Block
	}
	type args struct {
		loader    Loader
		chunkSize int
	}

	for _, data := range []struct {
		name          string
		fields        fields
		args          args
		wantLastBlock Block
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "success with leftDifficulty > rightDifficulty",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.On("Difficulty", "hash #3.2").Return(23, nil)
							proofer.On("Difficulty", "hash #3.1").Return(42, nil)
							proofer.On("Difficulty", "hash #3").Return(12, nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						blocks := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
						}

						storage := new(MockGroupStorage)
						storage.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data: func() Data {
								data := new(MockData)
								data.
									On("Equal", mock.AnythingOfType("*blockchain.MockData")).
									Return(true)

								return data
							}(),
							Hash:     "hash #2",
							PrevHash: "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return loader
				}(),
				chunkSize: 23,
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with leftDifficulty < rightDifficulty",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.On("Difficulty", "hash #3.2").Return(23, nil)
							proofer.On("Difficulty", "hash #3.1").Return(42, nil)
							proofer.On("Difficulty", "hash #3").Return(100, nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						blocks := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
						}
						blocksForDeleting := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
						}
						blocksForStoring := BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
						}
						newLastBlock := Block{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						}

						storage := new(MockGroupStorage)
						storage.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)
						storage.On("DeleteBlockGroup", blocksForDeleting).Return(nil)
						storage.On("StoreBlockGroup", blocksForStoring).Return(nil)
						storage.On("LoadLastBlock").Return(newLastBlock, nil)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data: func() Data {
								data := new(MockData)
								data.
									On("Equal", mock.AnythingOfType("*blockchain.MockData")).
									Return(true)

								return data
							}(),
							Hash:     "hash #2",
							PrevHash: "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return loader
				}(),
				chunkSize: 23,
			},
			wantLastBlock: Block{
				Timestamp: clock().Add(2 * time.Hour),
				Data:      new(MockData),
				Hash:      "hash #3",
				PrevHash:  "hash #2",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error with searching differences",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Proofer: new(MockProofer),
					},
					Storage: func() GroupStorage {
						storage := new(MockGroupStorage)
						storage.On("LoadBlocks", nil, 23).Return(nil, nil, iotest.ErrTimeout)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				loader:    new(MockLoader),
				chunkSize: 23,
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.Error,
		},
		{
			name: "error with calculating the difficulty of the left differences",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.On("Difficulty", "hash #3.2").Return(0, iotest.ErrTimeout)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						blocks := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
						}

						storage := new(MockGroupStorage)
						storage.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data: func() Data {
								data := new(MockData)
								data.
									On("Equal", mock.AnythingOfType("*blockchain.MockData")).
									Return(true)

								return data
							}(),
							Hash:     "hash #2",
							PrevHash: "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return loader
				}(),
				chunkSize: 23,
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.Error,
		},
		{
			name: "error with calculating the difficulty of the right differences",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.On("Difficulty", "hash #3.2").Return(23, nil)
							proofer.On("Difficulty", "hash #3.1").Return(42, nil)
							proofer.On("Difficulty", "hash #3").Return(0, iotest.ErrTimeout)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						blocks := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
						}

						storage := new(MockGroupStorage)
						storage.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data: func() Data {
								data := new(MockData)
								data.
									On("Equal", mock.AnythingOfType("*blockchain.MockData")).
									Return(true)

								return data
							}(),
							Hash:     "hash #2",
							PrevHash: "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return loader
				}(),
				chunkSize: 23,
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.Error,
		},
		{
			name: "error with leftDifficulty == rightDifficulty",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.On("Difficulty", "hash #3.2").Return(23, nil)
							proofer.On("Difficulty", "hash #3.1").Return(42, nil)
							proofer.On("Difficulty", "hash #3").Return(65, nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						blocks := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
						}

						storage := new(MockGroupStorage)
						storage.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data: func() Data {
								data := new(MockData)
								data.
									On("Equal", mock.AnythingOfType("*blockchain.MockData")).
									Return(true)

								return data
							}(),
							Hash:     "hash #2",
							PrevHash: "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return loader
				}(),
				chunkSize: 23,
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: func(
				test assert.TestingT,
				err error,
				msgAndArgs ...interface{},
			) bool {
				return assert.Equal(test, ErrEqualDifficulties, err, msgAndArgs...)
			},
		},
		{
			name: "error with deleting the left differences",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.On("Difficulty", "hash #3.2").Return(23, nil)
							proofer.On("Difficulty", "hash #3.1").Return(42, nil)
							proofer.On("Difficulty", "hash #3").Return(100, nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						blocks := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
						}
						blocksForDeleting := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
						}

						storage := new(MockGroupStorage)
						storage.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)
						storage.
							On("DeleteBlockGroup", blocksForDeleting).
							Return(iotest.ErrTimeout)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data: func() Data {
								data := new(MockData)
								data.
									On("Equal", mock.AnythingOfType("*blockchain.MockData")).
									Return(true)

								return data
							}(),
							Hash:     "hash #2",
							PrevHash: "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return loader
				}(),
				chunkSize: 23,
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.Error,
		},
		{
			name: "error with storing the right differences",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.On("Difficulty", "hash #3.2").Return(23, nil)
							proofer.On("Difficulty", "hash #3.1").Return(42, nil)
							proofer.On("Difficulty", "hash #3").Return(100, nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						blocks := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
						}
						blocksForDeleting := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
						}
						blocksForStoring := BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
						}

						storage := new(MockGroupStorage)
						storage.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)
						storage.On("DeleteBlockGroup", blocksForDeleting).Return(nil)
						storage.On("StoreBlockGroup", blocksForStoring).Return(iotest.ErrTimeout)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data: func() Data {
								data := new(MockData)
								data.
									On("Equal", mock.AnythingOfType("*blockchain.MockData")).
									Return(true)

								return data
							}(),
							Hash:     "hash #2",
							PrevHash: "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return loader
				}(),
				chunkSize: 23,
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.Error,
		},
		{
			name: "error with loading the last block",
			fields: fields{
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.On("Difficulty", "hash #3.2").Return(23, nil)
							proofer.On("Difficulty", "hash #3.1").Return(42, nil)
							proofer.On("Difficulty", "hash #3").Return(100, nil)

							return proofer
						}(),
					},
					Storage: func() GroupStorage {
						blocks := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
						}
						blocksForDeleting := BlockGroup{
							{
								Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.2",
								PrevHash:  "hash #3.1",
							},
							{
								Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
								Data:      new(MockData),
								Hash:      "hash #3.1",
								PrevHash:  "hash #2",
							},
						}
						blocksForStoring := BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
						}

						storage := new(MockGroupStorage)
						storage.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)
						storage.On("DeleteBlockGroup", blocksForDeleting).Return(nil)
						storage.On("StoreBlockGroup", blocksForStoring).Return(nil)
						storage.On("LoadLastBlock").Return(Block{}, iotest.ErrTimeout)

						return storage
					}(),
				},
				lastBlock: Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data: func() Data {
								data := new(MockData)
								data.
									On("Equal", mock.AnythingOfType("*blockchain.MockData")).
									Return(true)

								return data
							}(),
							Hash:     "hash #2",
							PrevHash: "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return loader
				}(),
				chunkSize: 23,
			},
			wantLastBlock: Block{
				Timestamp: clock(),
				Data:      new(MockData),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			blockchain := &Blockchain{
				dependencies: data.fields.dependencies,
				lastBlock:    data.fields.lastBlock,
			}
			gotErr := blockchain.Merge(data.args.loader, data.args.chunkSize)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.dependencies.Proofer,
				data.fields.dependencies.Storage,
				data.fields.lastBlock.Data,
				data.args.loader,
				blockchain.lastBlock.Data,
			)
			assert.Equal(test, data.wantLastBlock, blockchain.lastBlock)
			data.wantErr(test, gotErr)
		})
	}
}
