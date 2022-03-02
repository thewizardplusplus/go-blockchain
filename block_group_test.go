package blockchain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBlockGroup_IsValid(test *testing.T) {
	type args struct {
		prependedChunk BlockGroup
		validationMode ValidationMode
		proofer        Proofer
	}

	for _, data := range []struct {
		name   string
		blocks BlockGroup
		args   args
		want   assert.ErrorAssertionFunc
	}{
		{
			name:   "success without blocks",
			blocks: nil,
			args: args{
				prependedChunk: nil,
				validationMode: AsFullBlockchain,
				proofer:        new(MockProofer),
			},
			want: assert.NoError,
		},
		{
			name: "success with blocks as a full blockchain",
			blocks: BlockGroup{
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
					PrevHash:  "",
				},
			},
			args: args{
				prependedChunk: nil,
				validationMode: AsFullBlockchain,
				proofer: func() Proofer {
					proofer := new(MockProofer)
					for _, block := range (BlockGroup{
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
							PrevHash:  "",
						},
					}) {
						proofer.On("Validate", block).Return(nil)
					}

					return proofer
				}(),
			},
			want: assert.NoError,
		},
		{
			name: "success with blocks as a blockchain chunk",
			blocks: BlockGroup{
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
			args: args{
				prependedChunk: nil,
				validationMode: AsBlockchainChunk,
				proofer: func() Proofer {
					proofer := new(MockProofer)
					for _, block := range (BlockGroup{
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
					}) {
						proofer.On("Validate", block).Return(nil)
					}

					return proofer
				}(),
			},
			want: assert.NoError,
		},
		{
			name: "success with a prepended chunk",
			blocks: BlockGroup{
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
			},
			args: args{
				prependedChunk: BlockGroup{
					{
						Timestamp: clock().Add(3 * time.Hour),
						Data:      new(MockData),
						Hash:      "hash #4",
						PrevHash:  "hash #3",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockData),
						Hash:      "hash #3",
						PrevHash:  "hash #2",
					},
				},
				validationMode: AsFullBlockchain,
				proofer: func() Proofer {
					proofer := new(MockProofer)
					for _, block := range (BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
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
					}) {
						proofer.On("Validate", block).Return(nil)
					}

					return proofer
				}(),
			},
			want: assert.NoError,
		},
		{
			name: "failure due to the block in the middle",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockData),
					Hash:      "next incorrect hash",
					PrevHash:  "incorrect hash",
				},
				{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "",
				},
			},
			args: args{
				prependedChunk: nil,
				validationMode: AsFullBlockchain,
				proofer:        new(MockProofer),
			},
			want: assert.Error,
		},
		{
			name: "failure due to the block at the end (as a full blockchain)",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockData),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: time.Time{},
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "",
				},
			},
			args: args{
				prependedChunk: nil,
				validationMode: AsFullBlockchain,
				proofer: func() Proofer {
					proofer := new(MockProofer)
					proofer.
						On("Validate", Block{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "next hash",
							PrevHash:  "hash",
						}).
						Return(nil)

					return proofer
				}(),
			},
			want: assert.Error,
		},
		{
			name: "failure due to the block at the end (as a blockchain chunk)",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockData),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: time.Time{},
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				prependedChunk: nil,
				validationMode: AsBlockchainChunk,
				proofer: func() Proofer {
					proofer := new(MockProofer)
					proofer.
						On("Validate", Block{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "next hash",
							PrevHash:  "hash",
						}).
						Return(nil)

					return proofer
				}(),
			},
			want: assert.Error,
		},
		{
			name: "failure due to a prepended chunk",
			blocks: BlockGroup{
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
			},
			args: args{
				prependedChunk: BlockGroup{
					{
						Timestamp: clock().Add(3 * time.Hour),
						Data:      new(MockData),
						Hash:      "hash #4",
						PrevHash:  "hash #3",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockData),
						Hash:      "incorrect hash #2",
						PrevHash:  "incorrect hash #1",
					},
				},
				validationMode: AsFullBlockchain,
				proofer:        new(MockProofer),
			},
			want: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.blocks.IsValid(
				data.args.prependedChunk,
				data.args.validationMode,
				data.args.proofer,
			)

			for _, block := range data.args.prependedChunk {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			for _, block := range data.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			mock.AssertExpectationsForObjects(test, data.args.proofer)
			data.want(test, got)
		})
	}
}

func TestBlockGroup_IsLastBlockValid(test *testing.T) {
	type args struct {
		prevBlock      *Block
		validationMode ValidationMode
		proofer        Proofer
	}

	for _, data := range []struct {
		name   string
		blocks BlockGroup
		args   args
		want   assert.ErrorAssertionFunc
	}{
		{
			name: "success as a full blockchain",
			blocks: BlockGroup{
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
					PrevHash:  "",
				},
			},
			args: args{
				prevBlock:      nil,
				validationMode: AsFullBlockchain,
				proofer: func() Proofer {
					proofer := new(MockProofer)
					proofer.
						On("Validate", Block{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash",
							PrevHash:  "",
						}).
						Return(nil)

					return proofer
				}(),
			},
			want: assert.NoError,
		},
		{
			name: "success as a blockchain chunk without the previous block",
			blocks: BlockGroup{
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
			args: args{
				prevBlock:      nil,
				validationMode: AsBlockchainChunk,
				proofer: func() Proofer {
					proofer := new(MockProofer)
					proofer.
						On("Validate", Block{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash",
							PrevHash:  "previous hash",
						}).
						Return(nil)

					return proofer
				}(),
			},
			want: assert.NoError,
		},
		{
			name: "success as a blockchain chunk with the previous block",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockData),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			args: args{
				prevBlock: &Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash #2",
					PrevHash:  "hash #1",
				},
				validationMode: AsBlockchainChunk,
				proofer: func() Proofer {
					proofer := new(MockProofer)
					proofer.
						On("Validate", Block{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						}).
						Return(nil)

					return proofer
				}(),
			},
			want: assert.NoError,
		},
		{
			name: "failure as a full blockchain",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockData),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: time.Time{},
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "",
				},
			},
			args: args{
				prevBlock:      nil,
				validationMode: AsFullBlockchain,
				proofer:        new(MockProofer),
			},
			want: assert.Error,
		},
		{
			name: "failure as a blockchain chunk without the previous block",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockData),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: time.Time{},
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				prevBlock:      nil,
				validationMode: AsBlockchainChunk,
				proofer:        new(MockProofer),
			},
			want: assert.Error,
		},
		{
			name: "failure as a blockchain chunk with the previous block",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockData),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			args: args{
				prevBlock: &Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "incorrect hash #2",
					PrevHash:  "incorrect hash #1",
				},
				validationMode: AsBlockchainChunk,
				proofer:        new(MockProofer),
			},
			want: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.blocks.IsLastBlockValid(
				data.args.prevBlock,
				data.args.validationMode,
				data.args.proofer,
			)

			for _, block := range data.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			if data.args.prevBlock != nil {
				mock.AssertExpectationsForObjects(test, data.args.prevBlock.Data)
			}
			mock.AssertExpectationsForObjects(test, data.args.proofer)
			data.want(test, got)
		})
	}
}

func TestBlockGroup_FindBlock(test *testing.T) {
	type args struct {
		block Block
	}

	for _, data := range []struct {
		name           string
		blocks         BlockGroup
		args           args
		wantBlockIndex int
		wantIsFound    assert.BoolAssertionFunc
	}{
		{
			name: "is found",
			// the blocks should be sorted in descending order of their timestamps
			blocks: BlockGroup{
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
			},
			args: args{
				block: Block{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockData),
					Hash:      "hash #2",
					PrevHash:  "hash #1",
				},
			},
			wantBlockIndex: 1,
			wantIsFound:    assert.True,
		},
		{
			name: "is not found due to timestamp",
			// the blocks should be sorted in descending order of their timestamps
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #3",
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
			},
			args: args{
				block: Block{
					Timestamp: clock().Add(-time.Hour),
					Data:      new(MockData),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			wantBlockIndex: 0,
			wantIsFound:    assert.False,
		},
		{
			name: "is not found due to data",
			// the blocks should be sorted in descending order of their timestamps
			blocks: BlockGroup{
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
							Return(false)

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
			},
			args: args{
				block: Block{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockData),
					Hash:      "hash #2",
					PrevHash:  "hash #1",
				},
			},
			wantBlockIndex: 0,
			wantIsFound:    assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotBlockIndex, gotIsFound := data.blocks.FindBlock(data.args.block)

			for _, block := range data.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			mock.AssertExpectationsForObjects(test, data.args.block.Data)
			assert.Equal(test, data.wantBlockIndex, gotBlockIndex)
			data.wantIsFound(test, gotIsFound)
		})
	}
}

func TestBlockGroup_FindDifferences(test *testing.T) {
	type args struct {
		anotherBlocks BlockGroup
	}

	for _, data := range []struct {
		name           string
		blocks         BlockGroup
		args           args
		wantLeftIndex  int
		wantRightIndex int
		wantHasMatch   assert.BoolAssertionFunc
	}{
		{
			name: "same block groups",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #3",
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
			},
			args: args{
				anotherBlocks: BlockGroup{
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data: func() Data {
							data := new(MockData)
							data.
								On("Equal", mock.AnythingOfType("*blockchain.MockData")).
								Return(true)

							return data
						}(),
						Hash:     "hash #3",
						PrevHash: "hash #2",
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
				},
			},
			wantLeftIndex:  0,
			wantRightIndex: 0,
			wantHasMatch:   assert.True,
		},
		{
			name: "left block group is longer",
			blocks: BlockGroup{
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
			},
			args: args{
				anotherBlocks: BlockGroup{
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
				},
			},
			wantLeftIndex:  2,
			wantRightIndex: 1,
			wantHasMatch:   assert.True,
		},
		{
			name: "right block group is longer",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #3",
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
			},
			args: args{
				anotherBlocks: BlockGroup{
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
				},
			},
			wantLeftIndex:  1,
			wantRightIndex: 2,
			wantHasMatch:   assert.True,
		},
		{
			name: "different block groups",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #3",
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
			},
			args: args{
				anotherBlocks: BlockGroup{
					{
						Timestamp: clock().Add(2*time.Hour + 30*time.Minute),
						Data:      new(MockData),
						Hash:      "hash #3.1",
						PrevHash:  "hash #2.1",
					},
					{
						Timestamp: clock().Add(time.Hour + 30*time.Minute),
						Data:      new(MockData),
						Hash:      "hash #2.1",
						PrevHash:  "hash #1.1",
					},
					{
						Timestamp: clock().Add(30 * time.Minute),
						Data:      new(MockData),
						Hash:      "hash #1.1",
						PrevHash:  "",
					},
				},
			},
			wantLeftIndex:  0,
			wantRightIndex: 0,
			wantHasMatch:   assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLeftIndex, gotRightIndex, gotHasMatch :=
				data.blocks.FindDifferences(data.args.anotherBlocks)

			for _, block := range data.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			for _, block := range data.args.anotherBlocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			assert.Equal(test, data.wantLeftIndex, gotLeftIndex)
			assert.Equal(test, data.wantRightIndex, gotRightIndex)
			data.wantHasMatch(test, gotHasMatch)
		})
	}
}
