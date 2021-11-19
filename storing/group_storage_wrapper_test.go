package storing

import (
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestGroupStorageWrapper_StoreBlockGroup(test *testing.T) {
	type fields struct {
		Storage blockchain.Storage
	}
	type args struct {
		blocks blockchain.BlockGroup
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success without blocks",
			fields: fields{
				Storage: new(MockStorage),
			},
			args: args{
				blocks: nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with blocks",
			fields: fields{
				Storage: func() blockchain.Storage {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock(),
							Data:      new(MockStringer),
							Hash:      "hash #1",
							PrevHash:  "",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockStringer),
							Hash:      "hash #2",
							PrevHash:  "hash #1",
						},
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockStringer),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
					}

					storage := new(MockStorage)
					for _, block := range blocks {
						storage.On("StoreBlock", block).Return(nil)
					}

					return storage
				}(),
			},
			args: args{
				blocks: blockchain.BlockGroup{
					{
						Timestamp: clock(),
						Data:      new(MockStringer),
						Hash:      "hash #1",
						PrevHash:  "",
					},
					{
						Timestamp: clock().Add(time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #2",
						PrevHash:  "hash #1",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #3",
						PrevHash:  "hash #2",
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				Storage: func() blockchain.Storage {
					block := blockchain.Block{
						Timestamp: clock(),
						Data:      new(MockStringer),
						Hash:      "hash #1",
						PrevHash:  "",
					}

					storage := new(MockStorage)
					storage.On("StoreBlock", block).Return(iotest.ErrTimeout)

					return storage
				}(),
			},
			args: args{
				blocks: blockchain.BlockGroup{
					{
						Timestamp: clock(),
						Data:      new(MockStringer),
						Hash:      "hash #1",
						PrevHash:  "",
					},
					{
						Timestamp: clock().Add(time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #2",
						PrevHash:  "hash #1",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #3",
						PrevHash:  "hash #2",
					},
				},
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			wrapper := GroupStorageWrapper{
				Storage: data.fields.Storage,
			}
			gotErr := wrapper.StoreBlockGroup(data.args.blocks)

			mock.AssertExpectationsForObjects(test, data.fields.Storage)
			for _, block := range data.args.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			data.wantErr(test, gotErr)
		})
	}
}

func TestGroupStorageWrapper_DeleteBlockGroup(test *testing.T) {
	type fields struct {
		Storage blockchain.Storage
	}
	type args struct {
		blocks blockchain.BlockGroup
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success without blocks",
			fields: fields{
				Storage: new(MockStorage),
			},
			args: args{
				blocks: nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with blocks",
			fields: fields{
				Storage: func() blockchain.Storage {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock(),
							Data:      new(MockStringer),
							Hash:      "hash #1",
							PrevHash:  "",
						},
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockStringer),
							Hash:      "hash #2",
							PrevHash:  "hash #1",
						},
						{
							Timestamp: clock().Add(2 * time.Hour),
							Data:      new(MockStringer),
							Hash:      "hash #3",
							PrevHash:  "hash #2",
						},
					}

					storage := new(MockStorage)
					for _, block := range blocks {
						storage.On("DeleteBlock", block).Return(nil)
					}

					return storage
				}(),
			},
			args: args{
				blocks: blockchain.BlockGroup{
					{
						Timestamp: clock(),
						Data:      new(MockStringer),
						Hash:      "hash #1",
						PrevHash:  "",
					},
					{
						Timestamp: clock().Add(time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #2",
						PrevHash:  "hash #1",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #3",
						PrevHash:  "hash #2",
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				Storage: func() blockchain.Storage {
					block := blockchain.Block{
						Timestamp: clock(),
						Data:      new(MockStringer),
						Hash:      "hash #1",
						PrevHash:  "",
					}

					storage := new(MockStorage)
					storage.On("DeleteBlock", block).Return(iotest.ErrTimeout)

					return storage
				}(),
			},
			args: args{
				blocks: blockchain.BlockGroup{
					{
						Timestamp: clock(),
						Data:      new(MockStringer),
						Hash:      "hash #1",
						PrevHash:  "",
					},
					{
						Timestamp: clock().Add(time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #2",
						PrevHash:  "hash #1",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #3",
						PrevHash:  "hash #2",
					},
				},
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			wrapper := GroupStorageWrapper{
				Storage: data.fields.Storage,
			}
			gotErr := wrapper.DeleteBlockGroup(data.args.blocks)

			mock.AssertExpectationsForObjects(test, data.fields.Storage)
			for _, block := range data.args.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			data.wantErr(test, gotErr)
		})
	}
}

func clock() time.Time {
	year, month, day := 2006, time.January, 2
	hour, minute, second := 15, 4, 5
	return time.Date(
		year, month, day,
		hour, minute, second,
		0,        // nanosecond
		time.UTC, // location
	)
}
