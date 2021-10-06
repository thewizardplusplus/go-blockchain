package proofers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestProofOfWork_Hash(test *testing.T) {
	data := new(MockStringer)
	data.On("String").Return("hash")

	proofer := ProofOfWork{TargetBit: 248}
	hash := proofer.Hash(blockchain.Block{
		Timestamp: clock(),
		Data:      data,
		PrevHash:  "previous hash",
	})

	wantedHash :=
		"248:26:00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191"
	mock.AssertExpectationsForObjects(test, data)
	assert.Equal(test, wantedHash, hash)
}

func TestProofOfWork_Validate(test *testing.T) {
	type fields struct {
		TargetBit int
	}
	type args struct {
		block blockchain.Block
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   assert.BoolAssertionFunc
	}{
		{
			name:   "success",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() fmt.Stringer {
						data := new(MockStringer)
						data.On("String").Return("hash")

						return data
					}(),
					Hash: "248:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			want: assert.True,
		},
		{
			name:   "failure with another block data",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() fmt.Stringer {
						data := new(MockStringer)
						data.On("String").Return("hash #2")

						return data
					}(),
					Hash: "248:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			want: assert.False,
		},
		{
			name:   "failure with another nonce",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() fmt.Stringer {
						data := new(MockStringer)
						data.On("String").Return("hash")

						return data
					}(),
					Hash: "248:" +
						"42:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			want: assert.False,
		},
		{
			name:   "failure with another target bit",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() fmt.Stringer {
						data := new(MockStringer)
						data.On("String").Return("hash")

						return data
					}(),
					Hash: "42:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			want: assert.False,
		},
		{
			name:   "incorrect hash structure",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockStringer),
					Hash:      "incorrect",
					PrevHash:  "previous hash",
				},
			},
			want: assert.False,
		},
		{
			name:   "incorrect target bit",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockStringer),
					Hash: "incorrect:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			proofer := ProofOfWork{TargetBit: data.fields.TargetBit}
			got := proofer.Validate(data.args.block)

			mock.AssertExpectationsForObjects(test, data.args.block.Data)
			data.want(test, got)
		})
	}
}
