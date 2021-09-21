package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewMemoizingLoader(test *testing.T) {
	loader := new(MockLoader)
	memoizingLoader := NewMemoizingLoader(loader)

	mock.AssertExpectationsForObjects(test, loader)
	assert.Equal(test, loader, memoizingLoader.loader)
	assert.NotNil(test, memoizingLoader.loadingResults)
}
