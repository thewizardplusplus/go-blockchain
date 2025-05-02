package proofers

import (
	"bytes"
	"context"
	"math/big"
	"testing"
	"testing/iotest"
	"time"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
	powErrors "github.com/thewizardplusplus/go-pow/errors"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
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

func TestProofOfWork_HashEx(test *testing.T) {
	type fields struct {
		TargetBit                int
		MaxAttemptCount          mo.Option[int]
		RandomInitialNonceParams mo.Option[powValueTypes.RandomNonceParams]
	}
	type args struct {
		ctx   context.Context
		block blockchain.Block
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/zero initial nonce",
			fields: fields{
				TargetBit: 248,
			},
			args: args{
				ctx: context.Background(),
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					PrevHash: "previous hash",
				},
			},
			want: "248:" +
				"26:" +
				"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
			wantErr: assert.NoError,
		},
		{
			name: "success/random initial nonce",
			fields: fields{
				TargetBit: 248,
				RandomInitialNonceParams: mo.Some(powValueTypes.RandomNonceParams{
					RandomReader: bytes.NewReader([]byte("dummy")),
					MinRawValue:  big.NewInt(123),
					MaxRawValue:  big.NewInt(142),
				}),
			},
			args: args{
				ctx: context.Background(),
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					PrevHash: "previous hash",
				},
			},
			want: "248:" +
				"617:" +
				"00b68f5e223b82d5e3a4c1e48aca4db2b08f791866147ef33e3c5208937ae8f1",
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to construct the target bit index",
			fields: fields{
				TargetBit: -23,
			},
			args: args{
				ctx: context.Background(),
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockData),
					PrevHash:  "previous hash",
				},
			},
			want: "",
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, ErrInvalidParameters)
			},
		},
		{
			name: "error/unable to build the challenge",
			fields: fields{
				TargetBit: 1000,
			},
			args: args{
				ctx: context.Background(),
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					PrevHash: "previous hash",
				},
			},
			want: "",
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, ErrInvalidParameters)
			},
		},
		{
			name: "error/unable to solve the challenge/regular error",
			fields: fields{
				TargetBit: 248,
				RandomInitialNonceParams: mo.Some(powValueTypes.RandomNonceParams{
					RandomReader: bytes.NewReader([]byte("dummy")),
					MinRawValue:  big.NewInt(142),
					MaxRawValue:  big.NewInt(123),
				}),
			},
			args: args{
				ctx: context.Background(),
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					PrevHash: "previous hash",
				},
			},
			want: "",
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, ErrInvalidParameters)
			},
		},
		{
			name: "error/unable to solve the challenge/I/O error",
			fields: fields{
				TargetBit: 248,
				RandomInitialNonceParams: mo.Some(powValueTypes.RandomNonceParams{
					RandomReader: iotest.ErrReader(iotest.ErrTimeout),
					MinRawValue:  big.NewInt(123),
					MaxRawValue:  big.NewInt(142),
				}),
			},
			args: args{
				ctx: context.Background(),
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					PrevHash: "previous hash",
				},
			},
			want: "",
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, powErrors.ErrIO)
			},
		},
		{
			name: "error/unable to solve the challenge/context is done",
			fields: fields{
				TargetBit: 248,
			},
			args: args{
				ctx: func() context.Context {
					ctx, ctxCancel := context.WithCancel(context.Background())
					ctxCancel()

					return ctx
				}(),
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					PrevHash: "previous hash",
				},
			},
			want: "",
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, powErrors.ErrTaskInterruption)
			},
		},
		{
			name: "error/" +
				"unable to solve the challenge/" +
				"maximal attempt count is exceeded",
			fields: fields{
				TargetBit:       248,
				MaxAttemptCount: mo.Some(23),
			},
			args: args{
				ctx: context.Background(),
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					PrevHash: "previous hash",
				},
			},
			want: "",
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, powErrors.ErrTaskInterruption)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			proofer := ProofOfWork{
				TargetBit:                data.fields.TargetBit,
				MaxAttemptCount:          data.fields.MaxAttemptCount,
				RandomInitialNonceParams: data.fields.RandomInitialNonceParams,
			}
			got, err := proofer.HashEx(data.args.ctx, data.args.block)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)

			mock.AssertExpectationsForObjects(test, data.args.block.Data)
		})
	}
}

func TestProofOfWork_Validate(test *testing.T) {
	type args struct {
		block blockchain.Block
	}

	for _, data := range []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
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
			wantErr: assert.NoError,
		},
		{
			name: "error/the hash contains the invalid quantity of the parts",
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "invalid",
					PrevHash:  "previous hash",
				},
			},
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, ErrInvalidParameters)
			},
		},
		{
			name: "error/unable to parse the target bit",
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash: "invalid:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, ErrInvalidParameters)
			},
		},
		{
			name: "error/unable to construct the target bit index",
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash: "-23:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, ErrInvalidParameters)
			},
		},
		{
			name: "error/unable to parse the nonce",
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash: "248:" +
						"invalid:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, ErrInvalidParameters)
			},
		},
		{
			name: "error/unable to decode the hash sum",
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "248:26:invalid",
					PrevHash:  "previous hash",
				},
			},
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, ErrInvalidParameters)
			},
		},
		{
			name: "error/unable to build the challenge",
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					Hash: "1000:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, ErrInvalidParameters)
			},
		},
		{
			name: "error/unable to build the solution",
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("hash")

						return data
					}(),
					Hash:     "248:26:00c4c39529ced1cb3e32086b19b75383",
					PrevHash: "previous hash",
				},
			},
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, ErrInvalidParameters)
			},
		},
		{
			name: "error/unable to verify the solution",
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() blockchain.Data {
						data := new(MockData)
						data.On("String").Return("different hash")

						return data
					}(),
					Hash: "248:" +
						"26:" +
						"00c4c39529ced1cb3e32086b19b753831f6396c9fa79079bc93c1c76a6244191",
					PrevHash: "previous hash",
				},
			},
			wantErr: func(test assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(test, err, powErrors.ErrValidationFailure)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			err := (ProofOfWork{}).Validate(data.args.block)

			mock.AssertExpectationsForObjects(test, data.args.block.Data)
			data.wantErr(test, err)
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
