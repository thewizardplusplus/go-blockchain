package storages

import (
	"github.com/thewizardplusplus/go-blockchain"
)

// MemoryStorage ...
type MemoryStorage struct {
	blocks []blockchain.Block
}
