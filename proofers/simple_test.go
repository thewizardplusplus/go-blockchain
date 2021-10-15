package proofers

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestSimple_Hash(test *testing.T) {
	data := new(MockStringer)
	data.On("String").Return("hash")

	var proofer Simple
	hash := proofer.Hash(blockchain.Block{
		Timestamp: clock(),
		Data:      data,
		PrevHash:  "previous hash",
	})

	wantedHash :=
		"4a4292671f697950d1d1d3ec16967cacf0ca1c5e20a1e21b5e49712cf5e422ae"
	mock.AssertExpectationsForObjects(test, data)
	assert.Equal(test, wantedHash, hash)
}

func TestSimple_Validate(test *testing.T) {
	type args struct {
		block blockchain.Block
	}

	for _, data := range []struct {
		name string
		args args
		want assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() fmt.Stringer {
						data := new(MockStringer)
						data.On("String").Return("hash")

						return data
					}(),
					Hash: "4a4292671f697950d1d1d3ec16967cac" +
						"f0ca1c5e20a1e21b5e49712cf5e422ae",
					PrevHash: "previous hash",
				},
			},
			want: assert.NoError,
		},
		{
			name: "error",
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data: func() fmt.Stringer {
						data := new(MockStringer)
						data.On("String").Return("hash #2")

						return data
					}(),
					Hash: "4a4292671f697950d1d1d3ec16967cac" +
						"f0ca1c5e20a1e21b5e49712cf5e422ae",
					PrevHash: "previous hash",
				},
			},
			want: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var proofer Simple
			got := proofer.Validate(data.args.block)

			mock.AssertExpectationsForObjects(test, data.args.block.Data)
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
