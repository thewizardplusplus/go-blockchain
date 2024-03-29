package loading

import (
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestLoadStorage(test *testing.T) {
	type args struct {
		storage       blockchain.GroupStorage
		loader        blockchain.Loader
		initialCursor interface{}
		chunkSize     int
	}

	for _, data := range []struct {
		name           string
		args           args
		wantLastCursor interface{}
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "success without blocks",
			args: args{
				storage: new(MockGroupStorage),
				loader: func() blockchain.Loader {
					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(nil, "cursor-two", nil)

					return loader
				}(),
				initialCursor: "cursor-one",
				chunkSize:     23,
			},
			wantLastCursor: "cursor-one",
			wantErr:        assert.NoError,
		},
		{
			name: "success with blocks",
			args: args{
				storage: func() blockchain.GroupStorage {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "next hash",
							PrevHash:  "hash",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash",
							PrevHash:  "previous hash",
						},
					}

					storage := new(MockGroupStorage)
					storage.On("StoreBlockGroup", blocks).Return(nil)

					return storage
				}(),
				loader: func() blockchain.Loader {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "next hash",
							PrevHash:  "hash",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash",
							PrevHash:  "previous hash",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)
					loader.On("LoadBlocks", "cursor-two", 23).Return(nil, "cursor-three", nil)

					return loader
				}(),
				initialCursor: "cursor-one",
				chunkSize:     23,
			},
			wantLastCursor: "cursor-two",
			wantErr:        assert.NoError,
		},
		{
			name: "error with block loading",
			args: args{
				storage: new(MockGroupStorage),
				loader: func() blockchain.Loader {
					loader := new(MockLoader)
					loader.
						On("LoadBlocks", "cursor-one", 23).
						Return(nil, "", iotest.ErrTimeout)

					return loader
				}(),
				initialCursor: "cursor-one",
				chunkSize:     23,
			},
			wantLastCursor: "cursor-one",
			wantErr:        assert.Error,
		},
		{
			name: "error with block storing",
			args: args{
				storage: func() blockchain.GroupStorage {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "next hash",
							PrevHash:  "hash",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash",
							PrevHash:  "previous hash",
						},
					}

					storage := new(MockGroupStorage)
					storage.On("StoreBlockGroup", blocks).Return(iotest.ErrTimeout)

					return storage
				}(),
				loader: func() blockchain.Loader {
					blocks := blockchain.BlockGroup{
						{
							Timestamp: clock().Add(time.Hour),
							Data:      new(MockData),
							Hash:      "next hash",
							PrevHash:  "hash",
						},
						{
							Timestamp: clock(),
							Data:      new(MockData),
							Hash:      "hash",
							PrevHash:  "previous hash",
						},
					}

					loader := new(MockLoader)
					loader.On("LoadBlocks", "cursor-one", 23).Return(blocks, "cursor-two", nil)

					return loader
				}(),
				initialCursor: "cursor-one",
				chunkSize:     23,
			},
			wantLastCursor: "cursor-one",
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLastCursor, gotErr := LoadStorage(
				data.args.storage,
				data.args.loader,
				data.args.initialCursor,
				data.args.chunkSize,
			)

			mock.AssertExpectationsForObjects(test, data.args.storage, data.args.loader)
			assert.Equal(test, data.wantLastCursor, gotLastCursor)
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
