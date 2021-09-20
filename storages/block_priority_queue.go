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

// Less ...
func (queue BlockPriorityQueue) Less(i int, j int) bool {
	// use the descending order to the pop operation will return the last block
	return queue[i].Timestamp.After(queue[j].Timestamp)
}

// Swap ...
func (queue BlockPriorityQueue) Swap(i int, j int) {
	queue[i], queue[j] = queue[j], queue[i]
}
