package blockchain

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBlockchain(test *testing.T) {
	type args struct {
		genesisBlockData Hasher
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
				genesisBlockData: new(MockHasher),
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock:   clock,
						Proofer: new(MockProofer),
					},
					Storage: func() Storage {
						storage := new(MockStorage)
						storage.
							On("LoadLastBlock").
							Return(
								Block{
									Timestamp: clock(),
									Data:      new(MockHasher),
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
				Data:      new(MockHasher),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with an empty storage",
			args: args{
				genesisBlockData: new(MockHasher),
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock: clock,
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.
								On("Hash", Block{
									Timestamp: clock(),
									Data:      new(MockHasher),
									PrevHash:  "",
								}).
								Return("hash")

							return proofer
						}(),
					},
					Storage: func() Storage {
						storage := new(MockStorage)
						storage.On("LoadLastBlock").Return(Block{}, ErrEmptyStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockHasher),
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
				Data:      new(MockHasher),
				Hash:      "hash",
				PrevHash:  "",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error on last block loading",
			args: args{
				genesisBlockData: new(MockHasher),
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock:   clock,
						Proofer: new(MockProofer),
					},
					Storage: func() Storage {
						storage := new(MockStorage)
						storage.On("LoadLastBlock").Return(Block{}, iotest.ErrTimeout)

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
				genesisBlockData: new(MockHasher),
				dependencies: Dependencies{
					BlockDependencies: BlockDependencies{
						Clock: clock,
						Proofer: func() Proofer {
							proofer := new(MockProofer)
							proofer.
								On("Hash", Block{
									Timestamp: clock(),
									Data:      new(MockHasher),
									PrevHash:  "",
								}).
								Return("hash")

							return proofer
						}(),
					},
					Storage: func() Storage {
						storage := new(MockStorage)
						storage.On("LoadLastBlock").Return(Block{}, ErrEmptyStorage)
						storage.
							On("StoreBlock", Block{
								Timestamp: clock(),
								Data:      new(MockHasher),
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

func TestBlockchain_AddBlock(test *testing.T) {
	type fields struct {
		dependencies Dependencies
		lastBlock    Block
	}
	type args struct {
		data Hasher
	}

	for _, data := range []struct {
		name          string
		fields        fields
		args          args
		wantLastBlock Block
		wantErr       assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
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
				blockchain.lastBlock.Data,
			)
			assert.Equal(test, data.wantLastBlock, blockchain.lastBlock)
			data.wantErr(test, gotErr)
		})
	}
}
