package storages

import (
	"github.com/thewizardplusplus/go-blockchain"
)

// BlockPriorityQueue ...
type BlockPriorityQueue blockchain.BlockGroup

// Len ...
func (queue BlockPriorityQueue) Len() int {
	return len(queue)
}
