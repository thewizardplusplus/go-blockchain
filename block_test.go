package blockchain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBlock(test *testing.T) {
	data := new(MockHasher)
	proofer := new(MockProofer)
	proofer.
		On("Hash", Block{
			Timestamp: clock(),
			Data:      data,
			PrevHash:  "previous hash",
		}).
		Return("hash")

	prevBlock := Block{Hash: "previous hash"}
	block := NewBlock(data, prevBlock, Dependencies{
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
	data := new(MockHasher)
	proofer := new(MockProofer)
	proofer.
		On("Hash", Block{
			Timestamp: clock(),
			Data:      data,
			PrevHash:  "",
		}).
		Return("hash")

	block := NewGenesisBlock(data, Dependencies{
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

func TestBlock_IsValid(test *testing.T) {
	type fields struct {
		Timestamp time.Time
		Data      Hasher
		Hash      string
		PrevHash  string
	}
	type args struct {
		prevBlock    Block
		dependencies Dependencies
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   assert.BoolAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockHasher),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			args: args{
				prevBlock: Block{
					Timestamp: clock().Add(-time.Hour),
					Hash:      "previous hash",
				},
				dependencies: Dependencies{
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
			name: "failure due to timestamps",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockHasher),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			args: args{
				prevBlock: Block{
					Timestamp: clock().Add(time.Hour),
					Hash:      "previous hash",
				},
				dependencies: Dependencies{
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
			want: assert.False,
		},
		{
			name: "failure due to hashes",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockHasher),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			args: args{
				prevBlock: Block{
					Timestamp: clock().Add(-time.Hour),
					Hash:      "incorrect hash",
				},
				dependencies: Dependencies{
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
			want: assert.False,
		},
		{
			name: "failure due to proofers",
			fields: fields{
				Timestamp: clock(),
				Data:      new(MockHasher),
				Hash:      "hash",
				PrevHash:  "previous hash",
			},
			args: args{
				prevBlock: Block{
					Timestamp: clock().Add(-time.Hour),
					Hash:      "previous hash",
				},
				dependencies: Dependencies{
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
							Return(false)

						return proofer
					}(),
				},
			},
			want: assert.False,
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

			mock.AssertExpectationsForObjects(test, data.fields.Data)
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
