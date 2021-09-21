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
		dependencies   BlockDependencies
	}

	for _, data := range []struct {
		name   string
		blocks BlockGroup
		args   args
		want   assert.BoolAssertionFunc
	}{
		{
			name:   "success without blocks",
			blocks: nil,
			args: args{
				prependedChunk: nil,
				validationMode: AsFullBlockchain,
				dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			want: assert.True,
		},
		{
			name: "success with blocks as a full blockchain",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash",
					PrevHash:  "",
				},
			},
			args: args{
				prependedChunk: nil,
				validationMode: AsFullBlockchain,
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						for _, block := range (BlockGroup{
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockHasher),
								Hash:      "next hash",
								PrevHash:  "hash",
							},
							{
								Timestamp: clock(),
								Data:      new(MockHasher),
								Hash:      "hash",
								PrevHash:  "",
							},
						}) {
							proofer.On("Validate", block).Return(true)
						}

						return proofer
					}(),
				},
			},
			want: assert.True,
		},
		{
			name: "success with blocks as a blockchain chunk",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				prependedChunk: nil,
				validationMode: AsBlockchainChunk,
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						for _, block := range (BlockGroup{
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockHasher),
								Hash:      "next hash",
								PrevHash:  "hash",
							},
							{
								Timestamp: clock(),
								Data:      new(MockHasher),
								Hash:      "hash",
								PrevHash:  "previous hash",
							},
						}) {
							proofer.On("Validate", block).Return(true)
						}

						return proofer
					}(),
				},
			},
			want: assert.True,
		},
		{
			name: "success with a prepended chunk",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #2",
					PrevHash:  "hash #1",
				},
				{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash #1",
					PrevHash:  "",
				},
			},
			args: args{
				prependedChunk: BlockGroup{
					{
						Timestamp: clock().Add(3 * time.Hour),
						Data:      new(MockHasher),
						Hash:      "hash #4",
						PrevHash:  "hash #3",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockHasher),
						Hash:      "hash #3",
						PrevHash:  "hash #2",
					},
				},
				validationMode: AsFullBlockchain,
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						for _, block := range (BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockHasher),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockHasher),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
							{
								Timestamp: clock(),
								Data:      new(MockHasher),
								Hash:      "hash #1",
								PrevHash:  "",
							},
						}) {
							proofer.On("Validate", block).Return(true)
						}

						return proofer
					}(),
				},
			},
			want: assert.True,
		},
		{
			name: "failure due to the block in the middle",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "next incorrect hash",
					PrevHash:  "incorrect hash",
				},
				{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash",
					PrevHash:  "",
				},
			},
			args: args{
				prependedChunk: nil,
				validationMode: AsFullBlockchain,
				dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			want: assert.False,
		},
		{
			name: "failure due to the block at the end (as a full blockchain)",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: time.Time{},
					Data:      new(MockHasher),
					Hash:      "hash",
					PrevHash:  "",
				},
			},
			args: args{
				prependedChunk: nil,
				validationMode: AsFullBlockchain,
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockHasher),
								Hash:      "next hash",
								PrevHash:  "hash",
							}).
							Return(true)

						return proofer
					}(),
				},
			},
			want: assert.False,
		},
		{
			name: "failure due to the block at the end (as a blockchain chunk)",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: time.Time{},
					Data:      new(MockHasher),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				prependedChunk: nil,
				validationMode: AsBlockchainChunk,
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockHasher),
								Hash:      "next hash",
								PrevHash:  "hash",
							}).
							Return(true)

						return proofer
					}(),
				},
			},
			want: assert.False,
		},
		{
			name: "failure due to a prepended chunk",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #2",
					PrevHash:  "hash #1",
				},
				{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash #1",
					PrevHash:  "",
				},
			},
			args: args{
				prependedChunk: BlockGroup{
					{
						Timestamp: clock().Add(3 * time.Hour),
						Data:      new(MockHasher),
						Hash:      "hash #4",
						PrevHash:  "hash #3",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockHasher),
						Hash:      "incorrect hash #2",
						PrevHash:  "incorrect hash #1",
					},
				},
				validationMode: AsFullBlockchain,
				dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.blocks.IsValid(
				data.args.prependedChunk,
				data.args.validationMode,
				data.args.dependencies,
			)

			for _, block := range data.args.prependedChunk {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			for _, block := range data.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			mock.AssertExpectationsForObjects(test, data.args.dependencies.Proofer)
			data.want(test, got)
		})
	}
}

func TestBlockGroup_IsLastBlockValid(test *testing.T) {
	type args struct {
		prevBlock      *Block
		validationMode ValidationMode
		dependencies   BlockDependencies
	}

	for _, data := range []struct {
		name   string
		blocks BlockGroup
		args   args
		want   assert.BoolAssertionFunc
	}{
		{
			name: "success as a full blockchain",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash",
					PrevHash:  "",
				},
			},
			args: args{
				prevBlock:      nil,
				validationMode: AsFullBlockchain,
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock(),
								Data:      new(MockHasher),
								Hash:      "hash",
								PrevHash:  "",
							}).
							Return(true)

						return proofer
					}(),
				},
			},
			want: assert.True,
		},
		{
			name: "success as a blockchain chunk without the previous block",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				prevBlock:      nil,
				validationMode: AsBlockchainChunk,
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock(),
								Data:      new(MockHasher),
								Hash:      "hash",
								PrevHash:  "previous hash",
							}).
							Return(true)

						return proofer
					}(),
				},
			},
			want: assert.True,
		},
		{
			name: "success as a blockchain chunk with the previous block",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			args: args{
				prevBlock: &Block{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash #2",
					PrevHash:  "hash #1",
				},
				validationMode: AsBlockchainChunk,
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockHasher),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							}).
							Return(true)

						return proofer
					}(),
				},
			},
			want: assert.True,
		},
		{
			name: "failure as a full blockchain",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: time.Time{},
					Data:      new(MockHasher),
					Hash:      "hash",
					PrevHash:  "",
				},
			},
			args: args{
				prevBlock:      nil,
				validationMode: AsFullBlockchain,
				dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			want: assert.False,
		},
		{
			name: "failure as a blockchain chunk without the previous block",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: time.Time{},
					Data:      new(MockHasher),
					Hash:      "hash",
					PrevHash:  "previous hash",
				},
			},
			args: args{
				prevBlock:      nil,
				validationMode: AsBlockchainChunk,
				dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			want: assert.False,
		},
		{
			name: "failure as a blockchain chunk with the previous block",
			blocks: BlockGroup{
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			args: args{
				prevBlock: &Block{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "incorrect hash #2",
					PrevHash:  "incorrect hash #1",
				},
				validationMode: AsBlockchainChunk,
				dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.blocks.IsLastBlockValid(
				data.args.prevBlock,
				data.args.validationMode,
				data.args.dependencies,
			)

			for _, block := range data.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			if data.args.prevBlock != nil {
				mock.AssertExpectationsForObjects(test, data.args.prevBlock.Data)
			}
			mock.AssertExpectationsForObjects(test, data.args.dependencies.Proofer)
			data.want(test, got)
		})
	}
}
