package storages

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestMemoryStorage_LoadBlocks(test *testing.T) {
	type fields struct {
		blocks   blockchain.BlockGroup
		isSorted bool
	}
	type args struct {
		cursor interface{}
		count  int
	}

	for _, data := range []struct {
		name             string
		fields           fields
		args             args
		wantBlocks       blockchain.BlockGroup
		wantIsSorted     assert.BoolAssertionFunc
		wantLoadedBlocks blockchain.BlockGroup
		wantNextCursor   interface{}
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "with a nonempty storage and the unsorted blocks",
			fields: fields{
				blocks: blockchain.BlockGroup{
					{
						Timestamp: clock(),
						Data:      new(MockStringer),
						Hash:      "hash #1",
						PrevHash:  "",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #3",
						PrevHash:  "hash #2",
					},
					{
						Timestamp: clock().Add(time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #2",
						PrevHash:  "hash #1",
					},
				},
				isSorted: false,
			},
			args: args{
				cursor: 1,
				count:  2,
			},
			wantBlocks: blockchain.BlockGroup{
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
			wantIsSorted: assert.True,
			wantLoadedBlocks: blockchain.BlockGroup{
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
			wantNextCursor: 3,
			wantErr:        assert.NoError,
		},
		{
			name: "with a nonempty storage and the sorted blocks",
			fields: fields{
				blocks: blockchain.BlockGroup{
					{
						Timestamp: clock(),
						Data:      new(MockStringer),
						Hash:      "hash #1",
						PrevHash:  "",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #3",
						PrevHash:  "hash #2",
					},
					{
						Timestamp: clock().Add(time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #2",
						PrevHash:  "hash #1",
					},
				},
				isSorted: true,
			},
			args: args{
				cursor: 1,
				count:  2,
			},
			wantBlocks: blockchain.BlockGroup{
				{
					Timestamp: clock(),
					Data:      new(MockStringer),
					Hash:      "hash #1",
					PrevHash:  "",
				},
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #2",
					PrevHash:  "hash #1",
				},
			},
			wantIsSorted: assert.True,
			wantLoadedBlocks: blockchain.BlockGroup{
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
				{
					Timestamp: clock().Add(time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #2",
					PrevHash:  "hash #1",
				},
			},
			wantNextCursor: 3,
			wantErr:        assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			storage := MemoryStorage{
				blocks:   data.fields.blocks,
				isSorted: data.fields.isSorted,
			}
			gotLoadedBlocks, gotNextCursor, gotErr :=
				storage.LoadBlocks(data.args.cursor, data.args.count)

			for _, block := range storage.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			assert.Equal(test, data.wantBlocks, storage.blocks)
			data.wantIsSorted(test, storage.isSorted)
			assert.Equal(test, data.wantLoadedBlocks, gotLoadedBlocks)
			assert.Equal(test, data.wantNextCursor, gotNextCursor)
			data.wantErr(test, gotErr)
		})
	}
}

func TestMemoryStorage_LoadLastBlock(test *testing.T) {
	type fields struct {
		blocks    blockchain.BlockGroup
		lastBlock blockchain.Block
	}

	for _, data := range []struct {
		name          string
		fields        fields
		wantLastBlock blockchain.Block
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "success with the ascending order of the blocks",
			fields: fields{
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
				lastBlock: blockchain.Block{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			wantLastBlock: blockchain.Block{
				Timestamp: clock().Add(2 * time.Hour),
				Data:      new(MockStringer),
				Hash:      "hash #3",
				PrevHash:  "hash #2",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with the random order of the blocks",
			fields: fields{
				blocks: blockchain.BlockGroup{
					{
						Timestamp: clock(),
						Data:      new(MockStringer),
						Hash:      "hash #1",
						PrevHash:  "",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #3",
						PrevHash:  "hash #2",
					},
					{
						Timestamp: clock().Add(time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #2",
						PrevHash:  "hash #1",
					},
				},
				lastBlock: blockchain.Block{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			wantLastBlock: blockchain.Block{
				Timestamp: clock().Add(2 * time.Hour),
				Data:      new(MockStringer),
				Hash:      "hash #3",
				PrevHash:  "hash #2",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				blocks: nil,
			},
			wantLastBlock: blockchain.Block{},
			wantErr: func(
				test assert.TestingT,
				err error,
				msgAndArgs ...interface{},
			) bool {
				return assert.Equal(test, blockchain.ErrEmptyStorage, err, msgAndArgs...)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			storage := MemoryStorage{
				blocks:    data.fields.blocks,
				lastBlock: data.fields.lastBlock,
			}
			gotLastBlock, gotErr := storage.LoadLastBlock()

			for _, block := range data.fields.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			if data.fields.lastBlock != (blockchain.Block{}) {
				mock.AssertExpectationsForObjects(test, data.fields.lastBlock.Data)
			}
			assert.Equal(test, data.wantLastBlock, gotLastBlock)
			data.wantErr(test, gotErr)
		})
	}
}

func TestMemoryStorage_StoreBlock(test *testing.T) {
	type fields struct {
		blocks    blockchain.BlockGroup
		lastBlock blockchain.Block
		isSorted  bool
	}
	type args struct {
		block blockchain.Block
	}

	for _, data := range []struct {
		name          string
		fields        fields
		args          args
		wantBlocks    blockchain.BlockGroup
		wantLastBlock blockchain.Block
		wantIsSorted  assert.BoolAssertionFunc
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "success with a nonempty storage (with the adding of the last block)",
			fields: fields{
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
				lastBlock: blockchain.Block{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
				isSorted: true,
			},
			args: args{
				block: blockchain.Block{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
			},
			wantBlocks: blockchain.BlockGroup{
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
				{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
			},
			wantLastBlock: blockchain.Block{
				Timestamp: clock().Add(3 * time.Hour),
				Data:      new(MockStringer),
				Hash:      "hash #4",
				PrevHash:  "hash #3",
			},
			wantIsSorted: assert.False,
			wantErr:      assert.NoError,
		},
		{
			name: "success with a nonempty storage " +
				"(with the adding of the middle block)",
			fields: fields{
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
						Timestamp: clock().Add(3 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #4",
						PrevHash:  "hash #3",
					},
				},
				lastBlock: blockchain.Block{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				isSorted: true,
			},
			args: args{
				block: blockchain.Block{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			wantBlocks: blockchain.BlockGroup{
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
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			wantLastBlock: blockchain.Block{
				Timestamp: clock().Add(3 * time.Hour),
				Data:      new(MockStringer),
				Hash:      "hash #4",
				PrevHash:  "hash #3",
			},
			wantIsSorted: assert.False,
			wantErr:      assert.NoError,
		},
		{
			name: "success with an empty storage",
			fields: fields{
				blocks:    nil,
				lastBlock: blockchain.Block{},
				isSorted:  true,
			},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockStringer),
					Hash:      "hash #1",
					PrevHash:  "",
				},
			},
			wantBlocks: blockchain.BlockGroup{
				{
					Timestamp: clock(),
					Data:      new(MockStringer),
					Hash:      "hash #1",
					PrevHash:  "",
				},
			},
			wantLastBlock: blockchain.Block{
				Timestamp: clock(),
				Data:      new(MockStringer),
				Hash:      "hash #1",
				PrevHash:  "",
			},
			wantIsSorted: assert.False,
			wantErr:      assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			storage := &MemoryStorage{
				blocks:    data.fields.blocks,
				lastBlock: data.fields.lastBlock,
				isSorted:  data.fields.isSorted,
			}
			gotErr := storage.StoreBlock(data.args.block)

			for _, block := range data.fields.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			if data.fields.lastBlock != (blockchain.Block{}) {
				mock.AssertExpectationsForObjects(test, data.fields.lastBlock.Data)
			}
			mock.AssertExpectationsForObjects(test, data.args.block.Data)
			assert.Equal(test, data.wantBlocks, storage.blocks)
			assert.Equal(test, data.wantLastBlock, storage.lastBlock)
			data.wantIsSorted(test, storage.isSorted)
			data.wantErr(test, gotErr)
		})
	}
}

func TestMemoryStorage_DeleteBlock(test *testing.T) {
	type fields struct {
		blocks    blockchain.BlockGroup
		lastBlock blockchain.Block
		isSorted  bool
	}
	type args struct {
		block blockchain.Block
	}

	for _, data := range []struct {
		name          string
		fields        fields
		args          args
		wantBlocks    blockchain.BlockGroup
		wantLastBlock blockchain.Block
		wantIsSorted  assert.BoolAssertionFunc
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "with a nonempty storage, the deleting of the first block, " +
				"and the sorted blocks",
			fields: fields{
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
					{
						Timestamp: clock().Add(3 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #4",
						PrevHash:  "hash #3",
					},
				},
				lastBlock: blockchain.Block{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				isSorted: true,
			},
			args: args{
				block: blockchain.Block{
					Timestamp: clock(),
					Data:      new(MockStringer),
					Hash:      "hash #1",
					PrevHash:  "",
				},
			},
			wantBlocks: blockchain.BlockGroup{
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
				{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
			},
			wantLastBlock: blockchain.Block{
				Timestamp: clock().Add(3 * time.Hour),
				Data:      new(MockStringer),
				Hash:      "hash #4",
				PrevHash:  "hash #3",
			},
			wantIsSorted: assert.True,
			wantErr:      assert.NoError,
		},
		{
			name: "with a nonempty storage, the deleting of the middle block, " +
				"and the sorted blocks",
			fields: fields{
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
					{
						Timestamp: clock().Add(3 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #4",
						PrevHash:  "hash #3",
					},
				},
				lastBlock: blockchain.Block{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				isSorted: true,
			},
			args: args{
				block: blockchain.Block{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			wantBlocks: blockchain.BlockGroup{
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
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
			},
			wantLastBlock: blockchain.Block{
				Timestamp: clock().Add(3 * time.Hour),
				Data:      new(MockStringer),
				Hash:      "hash #4",
				PrevHash:  "hash #3",
			},
			wantIsSorted: assert.True,
			wantErr:      assert.NoError,
		},
		{
			name: "with a nonempty storage, the deleting of the middle block, " +
				"and the unsorted blocks",
			fields: fields{
				blocks: blockchain.BlockGroup{
					{
						Timestamp: clock(),
						Data:      new(MockStringer),
						Hash:      "hash #1",
						PrevHash:  "",
					},
					{
						Timestamp: clock().Add(3 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #4",
						PrevHash:  "hash #3",
					},
					{
						Timestamp: clock().Add(2 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #3",
						PrevHash:  "hash #2",
					},
					{
						Timestamp: clock().Add(time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #2",
						PrevHash:  "hash #1",
					},
				},
				lastBlock: blockchain.Block{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				isSorted: false,
			},
			args: args{
				block: blockchain.Block{
					Timestamp: clock().Add(2 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #3",
					PrevHash:  "hash #2",
				},
			},
			wantBlocks: blockchain.BlockGroup{
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
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
			},
			wantLastBlock: blockchain.Block{
				Timestamp: clock().Add(3 * time.Hour),
				Data:      new(MockStringer),
				Hash:      "hash #4",
				PrevHash:  "hash #3",
			},
			wantIsSorted: assert.True,
			wantErr:      assert.NoError,
		},
		{
			name: "with a nonempty storage, the deleting of the last block, " +
				"and the sorted blocks",
			fields: fields{
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
					{
						Timestamp: clock().Add(3 * time.Hour),
						Data:      new(MockStringer),
						Hash:      "hash #4",
						PrevHash:  "hash #3",
					},
				},
				lastBlock: blockchain.Block{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
				isSorted: true,
			},
			args: args{
				block: blockchain.Block{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
			},
			wantBlocks: blockchain.BlockGroup{
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
			wantLastBlock: blockchain.Block{
				Timestamp: clock().Add(2 * time.Hour),
				Data:      new(MockStringer),
				Hash:      "hash #3",
				PrevHash:  "hash #2",
			},
			wantIsSorted: assert.True,
			wantErr:      assert.NoError,
		},
		{
			name: "with an empty storage",
			fields: fields{
				blocks:    nil,
				lastBlock: blockchain.Block{},
				isSorted:  true,
			},
			args: args{
				block: blockchain.Block{
					Timestamp: clock().Add(3 * time.Hour),
					Data:      new(MockStringer),
					Hash:      "hash #4",
					PrevHash:  "hash #3",
				},
			},
			wantBlocks:    nil,
			wantLastBlock: blockchain.Block{},
			wantIsSorted:  assert.True,
			wantErr:       assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			storage := &MemoryStorage{
				blocks:    data.fields.blocks,
				lastBlock: data.fields.lastBlock,
				isSorted:  data.fields.isSorted,
			}
			gotErr := storage.DeleteBlock(data.args.block)

			for _, block := range data.fields.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			mock.AssertExpectationsForObjects(test, data.args.block.Data)
			assert.Equal(test, data.wantBlocks, storage.blocks)
			assert.Equal(test, data.wantLastBlock, storage.lastBlock)
			data.wantIsSorted(test, storage.isSorted)
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
