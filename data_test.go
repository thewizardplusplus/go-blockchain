package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewData(test *testing.T) {
	data := NewData("test")

	wantedData := universalDataWrapper{innerData: "test"}
	assert.Equal(test, wantedData, data)
}
