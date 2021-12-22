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
