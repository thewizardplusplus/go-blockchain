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

// Push ...
func (queue *BlockPriorityQueue) Push(block interface{}) {
	*queue = append(*queue, block.(blockchain.Block))
}

// Pop ...
func (queue *BlockPriorityQueue) Pop() interface{} {
	lastIndex := len(*queue) - 1
	lastBlock := (*queue)[lastIndex]
	// reset the popped block to avoid memory leaks via its reference fields
	(*queue)[lastIndex] = blockchain.Block{}

	*queue = (*queue)[:lastIndex]

	return lastBlock
}
