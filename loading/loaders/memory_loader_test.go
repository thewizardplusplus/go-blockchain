package loaders

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestMemoryLoader_LoadBlocks(test *testing.T) {
	type args struct {
		cursor interface{}
		count  int
	}

	for _, data := range []struct {
		name           string
		loader         MemoryLoader
		args           args
		wantBlocks     blockchain.BlockGroup
		wantNextCursor interface{}
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "with the loading of the chunk from the start",
			loader: MemoryLoader(blockchain.BlockGroup{
				{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
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
			}),
			args: args{
				cursor: nil,
				count:  2,
			},
			wantBlocks: blockchain.BlockGroup{
				{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			wantNextCursor: 2,
			wantErr:        assert.NoError,
		},
		{
			name: "with the loading of the chunk from the middle",
			loader: MemoryLoader(blockchain.BlockGroup{
				{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
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
			}),
			args: args{
				cursor: 1,
				count:  2,
			},
			wantBlocks: blockchain.BlockGroup{
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
			},
			wantNextCursor: 3,
			wantErr:        assert.NoError,
		},
		{
			name: "with the loading of the chunk from the end",
			loader: MemoryLoader(blockchain.BlockGroup{
				{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
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
			}),
			args: args{
				cursor: 3,
				count:  2,
			},
			wantBlocks: blockchain.BlockGroup{
				{
					Timestamp: clock(),
					Data:      new(MockData),
					Hash:      "hash #1",
					PrevHash:  "",
				},
			},
			wantNextCursor: 4,
			wantErr:        assert.NoError,
		},
		{
			name: "with the loading of the chunk from outside the range",
			loader: MemoryLoader(blockchain.BlockGroup{
				{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockData),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
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
			}),
			args: args{
				cursor: 10,
				count:  2,
			},
			wantBlocks:     blockchain.BlockGroup{},
			wantNextCursor: 4,
			wantErr:        assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotBlocks, gotNextCursor, gotErr :=
				data.loader.LoadBlocks(data.args.cursor, data.args.count)

			for _, block := range data.loader {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			for _, block := range data.wantBlocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			assert.Equal(test, data.wantBlocks, gotBlocks)
			assert.Equal(test, data.wantNextCursor, gotNextCursor)
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
