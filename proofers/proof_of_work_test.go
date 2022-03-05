package proofers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestProofOfWork_Hash(test *testing.T) {
	data := new(MockData)
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
		want   assert.ErrorAssertionFunc
	}{
		{
			name:   "success",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					Hash: "248:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			want: assert.NoError,
		},
		{
			name:   "failure with another block data",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash #2")

						return data
					}(),
					Hash: "248:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			want: assert.Error,
		},
		{
			name:   "failure with another nonce",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					Hash: "248:" +
						"42:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			want: assert.Error,
		},
		{
			name:   "failure with another target bit",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					Hash: "42:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			want: assert.Error,
		},
		{
			name:   "incorrect hash structure",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "incorrect",
					PrevHash:  "previous hash",
				},
			},
			want: assert.Error,
		},
		{
			name:   "incorrect target bit",
			fields: fields{TargetBit: 23},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash: "incorrect:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			want: assert.Error,
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

func TestProofOfWork_Difficulty(test *testing.T) {
	type fields struct {
		TargetBit int
	}
	type args struct {
		hash string
	}

	for _, data := range []struct {
		name           string
		fields         fields
		args           args
		wantDifficulty int
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name:   "success with target bit 248",
			fields: fields{TargetBit: 23},
			args: args{
				hash: "248:" +
					"26:" +
					"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
			},
			wantDifficulty: 7,
			wantErr:        assert.NoError,
		},
		{
			name:   "success with target bit 252",
			fields: fields{TargetBit: 23},
			args: args{
				hash: "252:" +
					"7:" +
					"0aef306f0c6e4d20e09a5b10bf3175b44f6cb19e836ad0a8bec96525bb2a9ef9",
			},
			wantDifficulty: 3,
			wantErr:        assert.NoError,
		},
		{
			name:   "incorrect hash structure",
			fields: fields{TargetBit: 23},
			args: args{
				hash: "incorrect",
			},
			wantDifficulty: 0,
			wantErr:        assert.Error,
		},
		{
			name:   "incorrect target bit",
			fields: fields{TargetBit: 23},
			args: args{
				hash: "incorrect:" +
					"26:" +
					"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
			},
			wantDifficulty: 0,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			proofer := ProofOfWork{TargetBit: data.fields.TargetBit}
			gotDifficulty, gotErr := proofer.Difficulty(data.args.hash)

			assert.Equal(test, data.wantDifficulty, gotDifficulty)
			data.wantErr(test, gotErr)
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
