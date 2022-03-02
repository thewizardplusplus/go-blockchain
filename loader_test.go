package blockchain

import (
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindDifferences(test *testing.T) {
	type args struct {
		leftLoader  Loader
		rightLoader Loader
		chunkSize   int
	}

	for _, data := range []struct {
		name                 string
		args                 args
		wantLeftDifferences  BlockGroup
		wantRightDifferences BlockGroup
		wantErr              assert.ErrorAssertionFunc
	}{
		{
			name: "success with same block groups",
			args: args{
				leftLoader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "hash #2",
							PrevHash:  "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					leftLoader := new(MockLoader)
					leftLoader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return leftLoader
				}(),
				rightLoader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data: func() Data {
								data := new(MockData)
								data.
									On("Equal", mock.AnythingOfType("*blockchain.MockData")).
									Return(true)

								return data
							}(),
							Hash:     "hash #3",
							PrevHash: "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "hash #2",
							PrevHash:  "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					rightLoader := new(MockLoader)
					rightLoader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return rightLoader
				}(),
				chunkSize: 23,
			},
			wantLeftDifferences:  BlockGroup{},
			wantRightDifferences: BlockGroup{},
			wantErr:              assert.NoError,
		},
		{
			name: "success with different block groups",
			args: args{
				leftLoader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
							Data:      new(MockData),
							Hash:      "hash #3.2",
							PrevHash:  "hash #3.1",
						},
						{
							Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
							Data:      new(MockData),
							Hash:      "hash #3.1",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "hash #2",
							PrevHash:  "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					leftLoader := new(MockLoader)
					leftLoader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return leftLoader
				}(),
				rightLoader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data: func() Data {
								data := new(MockData)
								data.
									On("Equal", mock.AnythingOfType("*blockchain.MockData")).
									Return(true)

								return data
							}(),
							Hash:     "hash #2",
							PrevHash: "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					rightLoader := new(MockLoader)
					rightLoader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return rightLoader
				}(),
				chunkSize: 23,
			},
			wantLeftDifferences: BlockGroup{
				{
					Timestamp: clock().Add(2*time.Hour + 40*time.Minute),
					Data:      new(MockData),
					Hash:      "hash #3.2",
					PrevHash:  "hash #3.1",
				},
				{
					Timestamp: clock().Add(2*time.Hour + 20*time.Minute),
					Data:      new(MockData),
					Hash:      "hash #3.1",
					PrevHash:  "hash #2",
				},
			},
			wantRightDifferences: BlockGroup{
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error with loading the left blocks",
			args: args{
				leftLoader: func() Loader {
					leftLoader := new(MockLoader)
					leftLoader.On("LoadBlocks", nil, 23).Return(nil, nil, iotest.ErrTimeout)

					return leftLoader
				}(),
				rightLoader: new(MockLoader),
				chunkSize:   23,
			},
			wantLeftDifferences:  nil,
			wantRightDifferences: nil,
			wantErr:              assert.Error,
		},
		{
			name: "error with loading the right blocks",
			args: args{
				leftLoader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "hash #2",
							PrevHash:  "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					leftLoader := new(MockLoader)
					leftLoader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return leftLoader
				}(),
				rightLoader: func() Loader {
					rightLoader := new(MockLoader)
					rightLoader.On("LoadBlocks", nil, 23).Return(nil, nil, iotest.ErrTimeout)

					return rightLoader
				}(),
				chunkSize: 23,
			},
			wantLeftDifferences:  nil,
			wantRightDifferences: nil,
			wantErr:              assert.Error,
		},
		{
			name: "error with different block groups",
			args: args{
				leftLoader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockData),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "hash #2",
							PrevHash:  "hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					leftLoader := new(MockLoader)
					leftLoader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return leftLoader
				}(),
				rightLoader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(2*time.Hour + 30*time.Minute),
							Data:      new(MockData),
							Hash:      "hash #3.1",
							PrevHash:  "hash #2.1",
						},
						{
							Timestamp: clock().Add(time.Hour + 30*time.Minute),
							Data:      new(MockData),
							Hash:      "hash #2.1",
							PrevHash:  "hash #1.1",
						},
						{
							Timestamp: clock().Add(30 * time.Minute),
							Data:      new(MockData),
							Hash:      "hash #1.1",
							PrevHash:  "",
						},
					}

					rightLoader := new(MockLoader)
					rightLoader.On("LoadBlocks", nil, 23).Return(blocks, 26, nil)

					return rightLoader
				}(),
				chunkSize: 23,
			},
			wantLeftDifferences:  nil,
			wantRightDifferences: nil,
			wantErr: func(
				test assert.TestingT,
				err error,
				msgAndArgs ...interface{},
			) bool {
				return assert.Equal(test, ErrNoMatch, err, msgAndArgs...)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLeftDifferences, gotRightDifferences, gotErr := FindDifferences(
				data.args.leftLoader,
				data.args.rightLoader,
				data.args.chunkSize,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.args.leftLoader,
				data.args.rightLoader,
			)
			assert.Equal(test, data.wantLeftDifferences, gotLeftDifferences)
			assert.Equal(test, data.wantRightDifferences, gotRightDifferences)
			data.wantErr(test, gotErr)
		})
	}
}
