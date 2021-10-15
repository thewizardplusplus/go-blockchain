package proofers

import (
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
	data := new(MockStringer)
	data.On("String").Return("hash")

	var proofer Simple
	err := proofer.Validate(blockchain.Block{
		Timestamp: clock(),
		Data:      data,
		Hash:      "4a4292671f697950d1d1d3ec16967cacf0ca1c5e20a1e21b5e49712cf5e422ae",
		PrevHash:  "previous hash",
	})

	mock.AssertExpectationsForObjects(test, data)
	assert.NoError(test, err)
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
