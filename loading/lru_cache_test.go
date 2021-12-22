package loading

import (
	"container/list"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestNewLRUCache(test *testing.T) {
	maximalSize := int(1e6)
	cache := NewLRUCache(maximalSize)

	assert.Equal(test, maximalSize, cache.maximalSize)
	assert.Equal(test, make(bucketGroup), cache.buckets)
	assert.Equal(test, list.New(), cache.queue)
}

func TestLRUCache_Get(test *testing.T) {
	type fields struct {
		buckets bucketGroup
		queue   *list.List
	}
	type args struct {
		parameters Parameters
	}

	for _, data := range []struct {
		name        string
		fields      fields
		args        args
		wantQueue   *list.List
		wantResults Results
		wantIsFound assert.BoolAssertionFunc
	}{
		{
			name: "with the existing element",
			fields: func() fields {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyTwo := Parameters{Cursor: "cursor #2", Count: 2}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(bucket{
					key: keyOne,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				})
				buckets[keyTwo] = queue.PushBack(bucket{
					key: keyTwo,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(3 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #4",
								PrevHash:  "hash #3",
							},
						},
						NextCursor: "cursor #3",
					},
				})

				return fields{
					buckets: buckets,
					queue:   queue,
				}
			}(),
			args: args{
				parameters: Parameters{Cursor: "cursor #2", Count: 2},
			},
			wantQueue: func() *list.List {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyTwo := Parameters{Cursor: "cursor #2", Count: 2}

				queue := list.New()
				queue.PushBack(bucket{
					key: keyTwo,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(3 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #4",
								PrevHash:  "hash #3",
							},
						},
						NextCursor: "cursor #3",
					},
				})
				queue.PushBack(bucket{
					key: keyOne,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				})

				return queue
			}(),
			wantResults: Results{
				Blocks: blockchain.BlockGroup{
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockData),
						Hash:      "hash #3",
						PrevHash:  "hash #2",
					},
					{
						Timestamp: clock().Add(3 * time.Hour),
						Data:      new(MockData),
						Hash:      "hash #4",
						PrevHash:  "hash #3",
					},
				},
				NextCursor: "cursor #3",
			},
			wantIsFound: assert.True,
		},
		{
			name: "with the non-existing element",
			fields: func() fields {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyTwo := Parameters{Cursor: "cursor #2", Count: 2}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(bucket{
					key: keyOne,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				})
				buckets[keyTwo] = queue.PushBack(bucket{
					key: keyTwo,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(3 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #4",
								PrevHash:  "hash #3",
							},
						},
						NextCursor: "cursor #3",
					},
				})

				return fields{
					buckets: buckets,
					queue:   queue,
				}
			}(),
			args: args{
				parameters: Parameters{Cursor: "cursor #3", Count: 2},
			},
			wantQueue: func() *list.List {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyTwo := Parameters{Cursor: "cursor #2", Count: 2}

				queue := list.New()
				queue.PushBack(bucket{
					key: keyOne,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				})
				queue.PushBack(bucket{
					key: keyTwo,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(3 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #4",
								PrevHash:  "hash #3",
							},
						},
						NextCursor: "cursor #3",
					},
				})

				return queue
			}(),
			wantResults: Results{},
			wantIsFound: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			cache := LRUCache{
				buckets: data.fields.buckets,
				queue:   data.fields.queue,
			}
			gotResults, gotIsFound := cache.Get(data.args.parameters)

			for _, bucketInstance := range data.fields.buckets {
				for _, block := range bucketInstance.Value.(bucket).value.Blocks {
					mock.AssertExpectationsForObjects(test, block.Data)
				}
			}
			assert.Equal(test, data.wantQueue, cache.queue)
			assert.Equal(test, data.wantResults, gotResults)
			data.wantIsFound(test, gotIsFound)
		})
	}
}

func TestLRUCache_Set(test *testing.T) {
	type fields struct {
		maximalSize int

		buckets bucketGroup
		queue   *list.List
	}
	type args struct {
		parameters Parameters
		results    Results
	}

	for _, data := range []struct {
		name       string
		fields     fields
		args       args
		wantFields fields
	}{
		{
			name: "with the existing element",
			fields: func() fields {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyTwo := Parameters{Cursor: "cursor #2", Count: 2}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(bucket{
					key: keyOne,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				})
				buckets[keyTwo] = queue.PushBack(bucket{
					key: keyTwo,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(3 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #4",
								PrevHash:  "hash #3",
							},
						},
						NextCursor: "cursor #3",
					},
				})

				return fields{
					maximalSize: 10,

					buckets: buckets,
					queue:   queue,
				}
			}(),
			args: args{
				parameters: Parameters{Cursor: "cursor #2", Count: 2},
				results: Results{
					Blocks: blockchain.BlockGroup{
						{
							Timestamp: clock().Add(4 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #5",
							PrevHash:  "hash #4",
						},
						{
							Timestamp: clock().Add(5 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #6",
							PrevHash:  "hash #5",
						},
					},
					NextCursor: "cursor #4",
				},
			},
			wantFields: func() fields {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyTwo := Parameters{Cursor: "cursor #2", Count: 2}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyTwo] = queue.PushBack(bucket{
					key: keyTwo,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(4 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #5",
								PrevHash:  "hash #4",
							},
							{
								Timestamp: clock().Add(5 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #6",
								PrevHash:  "hash #5",
							},
						},
						NextCursor: "cursor #4",
					},
				})
				buckets[keyOne] = queue.PushBack(bucket{
					key: keyOne,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				})

				return fields{
					buckets: buckets,
					queue:   queue,
				}
			}(),
		},
		{
			name: "with the non-existing element and the size less than the maximum",
			fields: func() fields {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyTwo := Parameters{Cursor: "cursor #2", Count: 2}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(bucket{
					key: keyOne,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				})
				buckets[keyTwo] = queue.PushBack(bucket{
					key: keyTwo,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(3 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #4",
								PrevHash:  "hash #3",
							},
						},
						NextCursor: "cursor #3",
					},
				})

				return fields{
					maximalSize: 10,

					buckets: buckets,
					queue:   queue,
				}
			}(),
			args: args{
				parameters: Parameters{Cursor: "cursor #3", Count: 2},
				results: Results{
					Blocks: blockchain.BlockGroup{
						{
							Timestamp: clock().Add(4 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #5",
							PrevHash:  "hash #4",
						},
						{
							Timestamp: clock().Add(5 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #6",
							PrevHash:  "hash #5",
						},
					},
					NextCursor: "cursor #4",
				},
			},
			wantFields: func() fields {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyTwo := Parameters{Cursor: "cursor #2", Count: 2}
				keyThree := Parameters{Cursor: "cursor #3", Count: 2}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyThree] = queue.PushBack(bucket{
					key: keyThree,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(4 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #5",
								PrevHash:  "hash #4",
							},
							{
								Timestamp: clock().Add(5 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #6",
								PrevHash:  "hash #5",
							},
						},
						NextCursor: "cursor #4",
					},
				})
				buckets[keyOne] = queue.PushBack(bucket{
					key: keyOne,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				})
				buckets[keyTwo] = queue.PushBack(bucket{
					key: keyTwo,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(3 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #4",
								PrevHash:  "hash #3",
							},
						},
						NextCursor: "cursor #3",
					},
				})

				return fields{
					buckets: buckets,
					queue:   queue,
				}
			}(),
		},
		{
			name: "with the non-existing element and the size greater than the maximum",
			fields: func() fields {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyTwo := Parameters{Cursor: "cursor #2", Count: 2}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(bucket{
					key: keyOne,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				})
				buckets[keyTwo] = queue.PushBack(bucket{
					key: keyTwo,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							},
							{
								Timestamp: clock().Add(3 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #4",
								PrevHash:  "hash #3",
							},
						},
						NextCursor: "cursor #3",
					},
				})

				return fields{
					maximalSize: 2,

					buckets: buckets,
					queue:   queue,
				}
			}(),
			args: args{
				parameters: Parameters{Cursor: "cursor #3", Count: 2},
				results: Results{
					Blocks: blockchain.BlockGroup{
						{
							Timestamp: clock().Add(4 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #5",
							PrevHash:  "hash #4",
						},
						{
							Timestamp: clock().Add(5 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #6",
							PrevHash:  "hash #5",
						},
					},
					NextCursor: "cursor #4",
				},
			},
			wantFields: func() fields {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyThree := Parameters{Cursor: "cursor #3", Count: 2}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyThree] = queue.PushBack(bucket{
					key: keyThree,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock().Add(4 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #5",
								PrevHash:  "hash #4",
							},
							{
								Timestamp: clock().Add(5 * time.Hour),
								Data:      new(MockData),
								Hash:      "hash #6",
								PrevHash:  "hash #5",
							},
						},
						NextCursor: "cursor #4",
					},
				})
				buckets[keyOne] = queue.PushBack(bucket{
					key: keyOne,
					value: Results{
						Blocks: blockchain.BlockGroup{
							{
								Timestamp: clock(),
								Data:      new(MockData),
								Hash:      "hash #1",
								PrevHash:  "",
							},
							{
								Timestamp: clock().Add(time.Hour),
								Data:      new(MockData),
								Hash:      "hash #2",
								PrevHash:  "hash #1",
							},
						},
						NextCursor: "cursor #2",
					},
				})

				return fields{
					buckets: buckets,
					queue:   queue,
				}
			}(),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			cache := LRUCache{
				maximalSize: data.fields.maximalSize,

				buckets: data.fields.buckets,
				queue:   data.fields.queue,
			}
			cache.Set(data.args.parameters, data.args.results)

			for _, bucketInstance := range data.fields.buckets {
				for _, block := range bucketInstance.Value.(bucket).value.Blocks {
					mock.AssertExpectationsForObjects(test, block.Data)
				}
			}
			for _, block := range data.args.results.Blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			assert.Equal(test, data.wantFields.buckets, cache.buckets)
			assert.Equal(test, data.wantFields.queue, cache.queue)
		})
	}
}

func TestLRUCache_getAndLiftElement(test *testing.T) {
	type fields struct {
		buckets     bucketGroup
		queue       *list.List
		wantElement *list.Element
	}
	type args struct {
		parameters Parameters
	}

	for _, data := range []struct {
		name        string
		fields      fields
		args        args
		wantQueue   *list.List
		wantIsFound assert.BoolAssertionFunc
	}{
		{
			name: "with the existing element",
			fields: func() fields {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyTwo := Parameters{Cursor: "cursor #2", Count: 2}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack("element #1")
				buckets[keyTwo] = queue.PushBack("element #2")

				return fields{
					buckets:     buckets,
					queue:       queue,
					wantElement: buckets[keyTwo],
				}
			}(),
			args: args{
				parameters: Parameters{Cursor: "cursor #2", Count: 2},
			},
			wantQueue: func() *list.List {
				queue := list.New()
				queue.PushBack("element #2")
				queue.PushBack("element #1")

				return queue
			}(),
			wantIsFound: assert.True,
		},
		{
			name: "with the non-existing element",
			fields: func() fields {
				keyOne := Parameters{Cursor: "cursor #1", Count: 2}
				keyTwo := Parameters{Cursor: "cursor #2", Count: 2}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack("element #1")
				buckets[keyTwo] = queue.PushBack("element #2")

				return fields{
					buckets:     buckets,
					queue:       queue,
					wantElement: nil,
				}
			}(),
			args: args{
				parameters: Parameters{Cursor: "cursor #3", Count: 2},
			},
			wantQueue: func() *list.List {
				queue := list.New()
				queue.PushBack("element #1")
				queue.PushBack("element #2")

				return queue
			}(),
			wantIsFound: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			cache := LRUCache{
				buckets: data.fields.buckets,
				queue:   data.fields.queue,
			}
			gotElement, gotIsFound := cache.getAndLiftElement(data.args.parameters)

			assert.Equal(test, data.wantQueue, cache.queue)
			assert.Equal(test, data.fields.wantElement, gotElement)
			data.wantIsFound(test, gotIsFound)
		})
	}
}
