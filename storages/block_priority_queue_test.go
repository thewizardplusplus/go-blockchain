package storages

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestBlockPriorityQueue_Len(test *testing.T) {
	queue := BlockPriorityQueue(blockchain.BlockGroup{
		{
			Timestamp: clock(),
			Data:      new(MockHasher),
			Hash:      "hash #1",
			PrevHash:  "",
		},
		{
			Timestamp: clock().Add(time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #2",
			PrevHash:  "hash #1",
		},
		{
			Timestamp: clock().Add(2 * time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #3",
			PrevHash:  "hash #2",
		},
	})
	gotLength := queue.Len()

	for _, block := range queue {
		mock.AssertExpectationsForObjects(test, block.Data)
	}
	assert.Equal(test, 3, gotLength)
}

func TestBlockPriorityQueue_Less(test *testing.T) {
	type args struct {
		i int
		j int
	}

	for _, data := range []struct {
		name       string
		queue      BlockPriorityQueue
		args       args
		wantIsLess assert.BoolAssertionFunc
	}{
		{
			name: "early",
			queue: BlockPriorityQueue(blockchain.BlockGroup{
				{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash #1",
					PrevHash:  "",
				},
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #2",
					PrevHash:  "hash #1",
				},
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			}),
			args: args{
				i: 0,
				j: 2,
			},
			wantIsLess: assert.False,
		},
		{
			name: "simultaneous",
			queue: BlockPriorityQueue(blockchain.BlockGroup{
				{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash #1",
					PrevHash:  "",
				},
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #2",
					PrevHash:  "hash #1",
				},
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			}),
			args: args{
				i: 1,
				j: 1,
			},
			wantIsLess: assert.False,
		},
		{
			name: "late",
			queue: BlockPriorityQueue(blockchain.BlockGroup{
				{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash #1",
					PrevHash:  "",
				},
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #2",
					PrevHash:  "hash #1",
				},
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			}),
			args: args{
				i: 2,
				j: 0,
			},
			wantIsLess: assert.True,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			isLess := data.queue.Less(data.args.i, data.args.j)

			for _, block := range data.queue {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			data.wantIsLess(test, isLess)
		})
	}
}

func TestBlockPriorityQueue_Swap(test *testing.T) {
	queue := BlockPriorityQueue(blockchain.BlockGroup{
		{
			Timestamp: clock(),
			Data:      new(MockHasher),
			Hash:      "hash #1",
			PrevHash:  "",
		},
		{
			Timestamp: clock().Add(time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #2",
			PrevHash:  "hash #1",
		},
		{
			Timestamp: clock().Add(2 * time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #3",
			PrevHash:  "hash #2",
		},
	})
	queue.Swap(0, 2)

	wantQueue := BlockPriorityQueue(blockchain.BlockGroup{
		{
			Timestamp: clock().Add(2 * time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #3",
			PrevHash:  "hash #2",
		},
		{
			Timestamp: clock().Add(time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #2",
			PrevHash:  "hash #1",
		},
		{
			Timestamp: clock(),
			Data:      new(MockHasher),
			Hash:      "hash #1",
			PrevHash:  "",
		},
	})
	for _, block := range queue {
		mock.AssertExpectationsForObjects(test, block.Data)
	}
	assert.Equal(test, wantQueue, queue)
}

func TestBlockPriorityQueue_Push(test *testing.T) {
	queue := BlockPriorityQueue(blockchain.BlockGroup{
		{
			Timestamp: clock(),
			Data:      new(MockHasher),
			Hash:      "hash #1",
			PrevHash:  "",
		},
		{
			Timestamp: clock().Add(time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #2",
			PrevHash:  "hash #1",
		},
		{
			Timestamp: clock().Add(2 * time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #3",
			PrevHash:  "hash #2",
		},
	})
	queue.Push(blockchain.Block{
		Timestamp: clock().Add(3 * time.Hour),
		Data:      new(MockHasher),
		Hash:      "hash #4",
		PrevHash:  "hash #3",
	})

	wantQueue := BlockPriorityQueue(blockchain.BlockGroup{
		{
			Timestamp: clock(),
			Data:      new(MockHasher),
			Hash:      "hash #1",
			PrevHash:  "",
		},
		{
			Timestamp: clock().Add(time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #2",
			PrevHash:  "hash #1",
		},
		{
			Timestamp: clock().Add(2 * time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #3",
			PrevHash:  "hash #2",
		},
		{
			Timestamp: clock().Add(3 * time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #4",
			PrevHash:  "hash #3",
		},
	})
	for _, block := range queue {
		mock.AssertExpectationsForObjects(test, block.Data)
	}
	assert.Equal(test, wantQueue, queue)
}

func TestBlockPriorityQueue_Pop(test *testing.T) {
	queue := BlockPriorityQueue(blockchain.BlockGroup{
		{
			Timestamp: clock(),
			Data:      new(MockHasher),
			Hash:      "hash #1",
			PrevHash:  "",
		},
		{
			Timestamp: clock().Add(time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #2",
			PrevHash:  "hash #1",
		},
		{
			Timestamp: clock().Add(2 * time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #3",
			PrevHash:  "hash #2",
		},
	})
	lastBlock := queue.Pop()

	wantQueue := BlockPriorityQueue(blockchain.BlockGroup{
		{
			Timestamp: clock(),
			Data:      new(MockHasher),
			Hash:      "hash #1",
			PrevHash:  "",
		},
		{
			Timestamp: clock().Add(time.Hour),
			Data:      new(MockHasher),
			Hash:      "hash #2",
			PrevHash:  "hash #1",
		},
	})
	for _, block := range queue {
		mock.AssertExpectationsForObjects(test, block.Data)
	}
	assert.Equal(test, wantQueue, queue)

	wantLastBlock := blockchain.Block{
		Timestamp: clock().Add(2 * time.Hour),
		Data:      new(MockHasher),
		Hash:      "hash #3",
		PrevHash:  "hash #2",
	}
	if assert.IsType(test, blockchain.Block{}, lastBlock) {
		mock.AssertExpectationsForObjects(test, lastBlock.(blockchain.Block).Data)
	}
	assert.Equal(test, wantLastBlock, lastBlock)
}
