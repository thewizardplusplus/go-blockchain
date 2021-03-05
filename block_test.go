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
