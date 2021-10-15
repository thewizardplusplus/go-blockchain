package blockchain

import (
	"fmt"
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBlock(test *testing.T) {
	data := new(MockStringer)
	proofer := new(MockProofer)
	proofer.
		On("Hash", Block{
			Timestamp: clock(),
			Data:      data,
			PrevHash:  "previous hash",
		}).
		Return("hash")

	prevBlock := Block{Hash: "previous hash"}
	block := NewBlock(data, prevBlock, BlockDependencies{
		Clock:   clock,
		Proofer: proofer,
	})

	wantedBlock := Block{
		Timestamp: clock(),
		Data:      data,
		Hash:      "hash",
		PrevHash:  "previous hash",
	}
	mock.AssertExpectationsForObjects(test, data, proofer)
	assert.Equal(test, wantedBlock, block)
}

func TestNewGenesisBlock(test *testing.T) {
	data := new(MockStringer)
	proofer := new(MockProofer)
	proofer.
		On("Hash", Block{
			Timestamp: clock(),
			Data:      data,
			PrevHash:  "",
		}).
		Return("hash")

	block := NewGenesisBlock(data, BlockDependencies{
		Clock:   clock,
		Proofer: proofer,
	})

	wantedBlock := Block{
		Timestamp: clock(),
		Data:      data,
		Hash:      "hash",
		PrevHash:  "",
	}
	mock.AssertExpectationsForObjects(test, data, proofer)
	assert.Equal(test, wantedBlock, block)
}

func TestBlock_MergedData(test *testing.T) {
	data := new(MockStringer)
	data.On("String").Return("hash")

	block := Block{
		Timestamp: clock(),
		Data:      data,
		Hash:      "hash",
		PrevHash:  "previous hash",
	}
	mergedData := block.MergedData()

	wantedMergedData := "2006-01-02 15:04:05 +0000 UTChashprevious hash"
	mock.AssertExpectationsForObjects(test, data)
	assert.Equal(test, wantedMergedData, mergedData)
}

func TestBlock_IsValid(test *testing.T) {
	type fields struct {
		Timestamp time.Time
		Data      fmt.Stringer
		Hash      string
		PrevHash  string
	}
	type args struct {
		prevBlock    *Block
		dependencies BlockDependencies
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   assert.ErrorAssertionFunc
	}{
		{
			name: "success with a previous block",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			args: args{
				prevBlock: &Block{
					Timestamp: clock().Add(-time.Hour),
					Hash:      "previous hash",
				},
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "hash",
								PrevHash:  "previous hash",
							}).
							Return(nil)

						return proofer
					}(),
				},
			},
			want: assert.NoError,
		},
		{
			name: "success without a previous block (with a previous hash)",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			args: args{
				prevBlock: nil,
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "hash",
								PrevHash:  "previous hash",
							}).
							Return(nil)

						return proofer
					}(),
				},
			},
			want: assert.NoError,
		},
		{
			name: "success without a previous block (without a previous hash)",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "",
			},
			args: args{
				prevBlock: nil,
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "hash",
								PrevHash:  "",
							}).
							Return(nil)

						return proofer
					}(),
				},
			},
			want: assert.NoError,
		},
		{
			name: "failure due to timestamps (with a previous block)",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			args: args{
				prevBlock: &Block{
					Timestamp: clock().Add(time.Hour),
					Hash:      "previous hash",
				},
				dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			want: assert.Error,
		},
		{
			name: "failure due to timestamps (without a previous block)",
			fields: fields{
				Timestamp: time.Time{},
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			args: args{
				prevBlock: nil,
				dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			want: assert.Error,
		},
		{
			name: "failure due to hashes",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			args: args{
				prevBlock: &Block{
					Timestamp: clock().Add(-time.Hour),
					Hash:      "incorrect hash",
				},
				dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			want: assert.Error,
		},
		{
			name: "failure due to proofers",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			args: args{
				prevBlock: &Block{
					Timestamp: clock().Add(-time.Hour),
					Hash:      "previous hash",
				},
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "hash",
								PrevHash:  "previous hash",
							}).
							Return(iotest.ErrTimeout)

						return proofer
					}(),
				},
			},
			want: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			block := Block{
				Timestamp: data.fields.Timestamp,
				Data:      data.fields.Data,
				Hash:      data.fields.Hash,
				PrevHash:  data.fields.PrevHash,
			}
			got := block.IsValid(data.args.prevBlock, data.args.dependencies)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.Data,
				data.args.dependencies.Proofer,
			)
			data.want(test, got)
		})
	}
}

func TestBlock_IsValidGenesisBlock(test *testing.T) {
	type fields struct {
		Timestamp time.Time
		Data      fmt.Stringer
		Hash      string
		PrevHash  string
	}
	type args struct {
		dependencies BlockDependencies
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "",
			},
			args: args{
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "hash",
								PrevHash:  "",
							}).
							Return(nil)

						return proofer
					}(),
				},
			},
			want: assert.NoError,
		},
		{
			name: "failure due to timestamps",
			fields: fields{
				Timestamp: time.Time{},
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "",
			},
			args: args{
				dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			want: assert.Error,
		},
		{
			name: "failure due to hashes",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			args: args{
				dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			want: assert.Error,
		},
		{
			name: "failure due to proofers",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash",
				PrevHash:  "",
			},
			args: args{
				dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock(),
								Data:      new(MockStringer),
								Hash:      "hash",
								PrevHash:  "",
							}).
							Return(iotest.ErrTimeout)

						return proofer
					}(),
				},
			},
			want: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			block := Block{
				Timestamp: data.fields.Timestamp,
				Data:      data.fields.Data,
				Hash:      data.fields.Hash,
				PrevHash:  data.fields.PrevHash,
			}
			got := block.IsValidGenesisBlock(data.args.dependencies)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.Data,
				data.args.dependencies.Proofer,
			)
			data.want(test, got)
		})
	}
}

func clock() time.Time {
	year, month, day := 2006, time.January, 2
	hour, minute, second := 15, 4, 5
	return time.Date(
		year, month, day,
		hour, minute, second,
		0,        // nanosecond
		time.UTC, // location
	)
}
