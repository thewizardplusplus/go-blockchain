package blockchain

import (
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLastBlockValidatingLoader_LoadBlocks(test *testing.T) {
	type fields struct {
		Loader       Loader
		Dependencies BlockDependencies
	}
	type args struct {
		cursor interface{}
		count  int
	}

	for _, data := range []struct {
		name           string
		fields         fields
		args           args
		wantBlocks     BlockGroup
		wantNextCursor interface{}
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "success without blocks",
			fields: fields{
				Loader: func() Loader {
					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(nil, "cursor-two", nil)

					return loader
				}(),
				Dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks:     nil,
			wantNextCursor: "cursor-two",
			wantErr:        assert.NoError,
		},
		{
			name: "success with blocks and without next blocks",
			fields: fields{
				Loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockHasher),
							Hash:      "next hash",
							PrevHash:  "hash",
						},
						{
							Timestamp: clock(),
							Data:      new(MockHasher),
							Hash:      "hash",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)
					loader.On("LoadBlocks", "cursor-two", 23).Return(nil, "cursor-three", nil)

					return loader
				}(),
				Dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock(),
								Data:      new(MockHasher),
								Hash:      "hash",
								PrevHash:  "",
							}).
							Return(true)

						return proofer
					}(),
				},
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks: BlockGroup{
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockHasher),
					Hash:      "next hash",
					PrevHash:  "hash",
				},
				{
					Timestamp: clock(),
					Data:      new(MockHasher),
					Hash:      "hash",
					PrevHash:  "",
				},
			},
			wantNextCursor: "cursor-two",
			wantErr:        assert.NoError,
		},
		{
			name: "success with blocks and next blocks",
			fields: fields{
				Loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(3 * time.Hour),
							Data:      new(MockHasher),
							Hash:      "hash #4",
							PrevHash:  "hash #3",
						},
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockHasher),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
					}
					nextBlocks := BlockGroup{
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
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)
					loader.
						On("LoadBlocks", "cursor-two", 23).
						Return(nextBlocks, "cursor-three", nil)

					return loader
				}(),
				Dependencies: BlockDependencies{
					Clock: clock,
					Proofer: func() Proofer {
						proofer := new(MockProofer)
						proofer.
							On("Validate", Block{
								Timestamp: clock().Add(2 * time.Hour),
								Data:      new(MockHasher),
								Hash:      "hash #3",
								PrevHash:  "hash #2",
							}).
							Return(true)

						return proofer
					}(),
				},
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks: BlockGroup{
				{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockHasher),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			wantNextCursor: "cursor-two",
			wantErr:        assert.NoError,
		},
		{
			name: "error with block loading",
			fields: fields{
				Loader: func() Loader {
					loader := new(MockLoader)
					loader.
						On("LoadBlocks", "cursor-one", 23).
						Return(nil, "", iotest.ErrTimeout)

					return loader
				}(),
				Dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks:     nil,
			wantNextCursor: nil,
			wantErr:        assert.Error,
		},
		{
			name: "error with next block loading",
			fields: fields{
				Loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockHasher),
							Hash:      "next hash",
							PrevHash:  "hash",
						},
						{
							Timestamp: clock(),
							Data:      new(MockHasher),
							Hash:      "hash",
							PrevHash:  "previous hash",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)
					loader.
						On("LoadBlocks", "cursor-two", 23).
						Return(nil, nil, iotest.ErrTimeout)

					return loader
				}(),
				Dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks:     nil,
			wantNextCursor: nil,
			wantErr:        assert.Error,
		},
		{
			name: "error with block validating and without next blocks",
			fields: fields{
				Loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockHasher),
							Hash:      "next hash",
							PrevHash:  "hash",
						},
						{
							Timestamp: time.Time{},
							Data:      new(MockHasher),
							Hash:      "hash",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)
					loader.On("LoadBlocks", "cursor-two", 23).Return(nil, "cursor-three", nil)

					return loader
				}(),
				Dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks:     nil,
			wantNextCursor: nil,
			wantErr:        assert.Error,
		},
		{
			name: "error with block validating and next blocks",
			fields: fields{
				Loader: func() Loader {
					blocks := BlockGroup{
						{
							Timestamp: clock().Add(3 * time.Hour),
							Data:      new(MockHasher),
							Hash:      "hash #4",
							PrevHash:  "hash #3",
						},
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockHasher),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
					}
					nextBlocks := BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockHasher),
							Hash:      "incorrect hash #2",
							PrevHash:  "incorrect hash #1",
						},
						{
							Timestamp: clock(),
							Data:      new(MockHasher),
							Hash:      "hash #1",
							PrevHash:  "",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)
					loader.
						On("LoadBlocks", "cursor-two", 23).
						Return(nextBlocks, "cursor-three", nil)

					return loader
				}(),
				Dependencies: BlockDependencies{
					Clock:   clock,
					Proofer: new(MockProofer),
				},
			},
			args: args{
				cursor: "cursor-one",
				count:  23,
			},
			wantBlocks:     nil,
			wantNextCursor: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			loader := LastBlockValidatingLoader{
				Loader:       data.fields.Loader,
				Dependencies: data.fields.Dependencies,
			}
			gotBlocks, gotNextCursor, gotErr :=
				loader.LoadBlocks(data.args.cursor, data.args.count)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.Loader,
				data.fields.Dependencies.Proofer,
			)
			assert.Equal(test, data.wantBlocks, gotBlocks)
			assert.Equal(test, data.wantNextCursor, gotNextCursor)
			data.wantErr(test, gotErr)
		})
	}
}
