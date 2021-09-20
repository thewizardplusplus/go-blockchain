package blockchain

import (
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewGroupStorage(test *testing.T) {
	type args struct {
		storage Storage
	}

	for _, data := range []struct {
		name             string
		args             args
		wantGroupStorage GroupStorage
	}{
		{
			name: "with the group storage",
			args: args{
				storage: new(MockGroupStorage),
			},
			wantGroupStorage: new(MockGroupStorage),
		},
		{
			name: "with the storage",
			args: args{
				storage: new(MockStorage),
			},
			wantGroupStorage: GroupStorageWrapper{Storage: new(MockStorage)},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotGroupStorage := NewGroupStorage(data.args.storage)

			mock.AssertExpectationsForObjects(test, data.args.storage)
			assert.Equal(test, data.wantGroupStorage, gotGroupStorage)
		})
	}
}

func TestGroupStorageWrapper_StoreBlockGroup(test *testing.T) {
	type fields struct {
		Storage Storage
	}
	type args struct {
		blocks BlockGroup
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
				Storage: func() Storage {
					blocks := BlockGroup{
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
					}

					storage := new(MockStorage)
					for _, block := range blocks {
						storage.On("StoreBlock", block).Return(nil)
					}

					return storage
				}(),
			},
			args: args{
				blocks: BlockGroup{
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
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				Storage: func() Storage {
					block := Block{
						Timestamp: clock(),
						Data:      new(MockHasher),
						Hash:      "hash #1",
						PrevHash:  "",
					}

					storage := new(MockStorage)
					storage.On("StoreBlock", block).Return(iotest.ErrTimeout)

					return storage
				}(),
			},
			args: args{
				blocks: BlockGroup{
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
