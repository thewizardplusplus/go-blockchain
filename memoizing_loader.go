package blockchain

import (
	"sync"
)

type loadingParameters struct {
	cursor interface{}
	count  int
}

type loadingResult struct {
	blocks     BlockGroup
	nextCursor interface{}
}

// MemoizingLoader ...
type MemoizingLoader struct {
	loader         Loader
	loadingResults *sync.Map
}
