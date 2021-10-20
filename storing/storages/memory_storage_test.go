package storages

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-blockchain"
)

func TestMemoryStorage_Blocks(test *testing.T) {
	storage := MemoryStorage{
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
	}
	gotBlocks := storage.Blocks()

	for _, block := range storage.blocks {
		mock.AssertExpectationsForObjects(test, block.Data)
	}
	assert.Equal(test, storage.blocks, gotBlocks)
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
	}
	type args struct {
		block blockchain.Block
	}

	for _, data := range []struct {
		name          string
		fields        fields
		args          args
		wantLastBlock blockchain.Block
		wantBlocks    blockchain.BlockGroup
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
			},
			args: args{
				block: blockchain.Block{
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
			wantErr: assert.NoError,
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
			},
			args: args{
				block: blockchain.Block{
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
			wantErr: assert.NoError,
		},
		{
			name: "success with an empty storage",
			fields: fields{
				blocks:    nil,
				lastBlock: blockchain.Block{},
			},
			args: args{
				block: blockchain.Block{
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
			wantBlocks: blockchain.BlockGroup{
				{
					Timestamp: clock(),
					Data:      new(MockStringer),
					Hash:      "hash #1",
					PrevHash:  "",
				},
			},
			wantErr: assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			storage := &MemoryStorage{
				blocks:    data.fields.blocks,
				lastBlock: data.fields.lastBlock,
			}
			gotErr := storage.StoreBlock(data.args.block)

			for _, block := range data.fields.blocks {
				mock.AssertExpectationsForObjects(test, block.Data)
			}
			if data.fields.lastBlock != (blockchain.Block{}) {
				mock.AssertExpectationsForObjects(test, data.fields.lastBlock.Data)
			}
			mock.AssertExpectationsForObjects(test, data.args.block.Data)
			assert.Equal(test, data.wantLastBlock, storage.lastBlock)
			assert.Equal(test, data.wantBlocks, storage.blocks)
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
