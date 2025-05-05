package blockchain_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/samber/mo"
	"github.com/thewizardplusplus/go-blockchain"
	"github.com/thewizardplusplus/go-blockchain/proofers"
	"github.com/thewizardplusplus/go-blockchain/storing"
	"github.com/thewizardplusplus/go-blockchain/storing/storages"
)

func ExampleBlockchain() {
	timestamp := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	blockDependencies := blockchain.BlockDependencies{
		// use the custom clock function to get the same blocks
		Clock: func() time.Time {
			timestamp = timestamp.Add(time.Hour)
			return timestamp
		},
		Proofer: proofers.ProofOfWork{TargetBit: 248},
	}

	blockchainInstance, err := blockchain.NewBlockchainEx(
		context.Background(),
		blockchain.NewBlockchainExParams{
			Clock:            blockDependencies.Clock,
			Proofer:          blockDependencies.Proofer,
			Storage:          storing.NewGroupStorage(&storages.MemoryStorage{}),
			GenesisBlockData: mo.Some(blockchain.NewData("genesis block")),
		},
	)
	if err != nil {
		log.Fatalf("unable to create a new blockchain: %v", err)
	}

	const blockCount = 5
	for i := 0; i < blockCount; i++ {
		if err := blockchainInstance.AddBlock(
			blockchain.NewData(fmt.Sprintf("block #%d", i)),
		); err != nil {
			log.Fatalf("unable to add the block: %v", err)
		}
	}

	addedBlocks, _, _ := blockchainInstance.LoadBlocks(nil, blockCount+1)
	blocksBytes, _ := json.MarshalIndent(addedBlocks, "", "  ")
	fmt.Println(string(blocksBytes))

	// Output:
	// [
	//   {
	//     "Timestamp": "2006-01-02T21:04:05Z",
	//     "Data": "block #4",
	//     "Hash": "248:173:00b6863763acd6ec77ca3521589d8e68c118efe855657d702783e8e6aee169a9",
	//     "PrevHash": "248:65:00d5800e119abe44d89469c2161be7f9645d7237697c6d14b4a72717893582fa"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T20:04:05Z",
	//     "Data": "block #3",
	//     "Hash": "248:65:00d5800e119abe44d89469c2161be7f9645d7237697c6d14b4a72717893582fa",
	//     "PrevHash": "248:136:003c7def3d467a759fad481c03cadbd62e62b2c5dbc10e4bbb6e1944c158a8be"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T19:04:05Z",
	//     "Data": "block #2",
	//     "Hash": "248:136:003c7def3d467a759fad481c03cadbd62e62b2c5dbc10e4bbb6e1944c158a8be",
	//     "PrevHash": "248:15:002fc891ad012c4a89f7b267a2ec1767415c627ff69b88b90a93be938b026efa"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T18:04:05Z",
	//     "Data": "block #1",
	//     "Hash": "248:15:002fc891ad012c4a89f7b267a2ec1767415c627ff69b88b90a93be938b026efa",
	//     "PrevHash": "248:198:0058f5dae6ca3451801a276c94862c7cce085e6f9371e50d80ddbb87c1438faf"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T17:04:05Z",
	//     "Data": "block #0",
	//     "Hash": "248:198:0058f5dae6ca3451801a276c94862c7cce085e6f9371e50d80ddbb87c1438faf",
	//     "PrevHash": "248:225:00e26abd9974fcdea4b32eca43c9dc5c67fffa8efd53cebffa9b049fd6c2bb36"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T16:04:05Z",
	//     "Data": "genesis block",
	//     "Hash": "248:225:00e26abd9974fcdea4b32eca43c9dc5c67fffa8efd53cebffa9b049fd6c2bb36",
	//     "PrevHash": ""
	//   }
	// ]
}

func ExampleBlockchain_Merge() {
	timestamp := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	blockGroupOne := blockchain.BlockGroup{
		{
			Timestamp: timestamp.Add(2*time.Hour + 40*time.Minute),
			Data:      blockchain.NewData("block #1.2"),
			Hash:      "250:0:hash #3.2",
			PrevHash:  "23:0:hash #3.1",
		},
		{
			Timestamp: timestamp.Add(2*time.Hour + 20*time.Minute),
			Data:      blockchain.NewData("block #1.1"),
			Hash:      "250:0:hash #3.1",
			PrevHash:  "23:0:hash #2",
		},
		{
			Timestamp: timestamp.Add(time.Hour),
			Data:      blockchain.NewData("block #0"),
			Hash:      "23:0:hash #2",
			PrevHash:  "23:0:hash #1",
		},
		{
			Timestamp: timestamp,
			Data:      blockchain.NewData("genesis block"),
			Hash:      "23:0:hash #1",
			PrevHash:  "",
		},
	}
	blockchainInstanceOne, err := blockchain.NewBlockchainEx(
		context.Background(),
		blockchain.NewBlockchainExParams{
			Clock: time.Now,
			Proofer: proofers.ProofOfWork{
				TargetBit: 248,
			},
			Storage: storing.NewGroupStorage(
				storages.NewMemoryStorage(blockGroupOne),
			),
			GenesisBlockData: mo.None[blockchain.Data](),
		},
	)
	if err != nil {
		log.Fatalf("unable to create the blockchain #1: %v", err)
	}

	blockGroupTwo := blockchain.BlockGroup{
		{
			Timestamp: timestamp.Add(2 * time.Hour),
			Data:      blockchain.NewData("block #1"),
			Hash:      "23:0:hash #3",
			PrevHash:  "23:0:hash #2",
		},
		{
			Timestamp: timestamp.Add(time.Hour),
			Data:      blockchain.NewData("block #0"),
			Hash:      "23:0:hash #2",
			PrevHash:  "23:0:hash #1",
		},
		{
			Timestamp: timestamp,
			Data:      blockchain.NewData("genesis block"),
			Hash:      "23:0:hash #1",
			PrevHash:  "",
		},
	}
	blockchainInstanceTwo, err := blockchain.NewBlockchainEx(
		context.Background(),
		blockchain.NewBlockchainExParams{
			Clock: time.Now,
			Proofer: proofers.ProofOfWork{
				TargetBit: 248,
			},
			Storage: storing.NewGroupStorage(
				storages.NewMemoryStorage(blockGroupTwo),
			),
			GenesisBlockData: mo.None[blockchain.Data](),
		},
	)
	if err != nil {
		log.Fatalf("unable to create the blockchain #2: %v", err)
	}

	if err := blockchainInstanceOne.Merge(blockchainInstanceTwo, 3); err != nil {
		log.Fatalf("unable to merge the blockchains: %v", err)
	}

	mergedBlocks, _, _ := blockchainInstanceOne.LoadBlocks(nil, 10)
	blocksBytes, _ := json.MarshalIndent(mergedBlocks, "", "  ")
	fmt.Println(string(blocksBytes))

	// Output:
	// [
	//   {
	//     "Timestamp": "2006-01-02T17:04:05Z",
	//     "Data": "block #1",
	//     "Hash": "23:0:hash #3",
	//     "PrevHash": "23:0:hash #2"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T16:04:05Z",
	//     "Data": "block #0",
	//     "Hash": "23:0:hash #2",
	//     "PrevHash": "23:0:hash #1"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T15:04:05Z",
	//     "Data": "genesis block",
	//     "Hash": "23:0:hash #1",
	//     "PrevHash": ""
	//   }
	// ]
}

func ExampleBlock() {
	timestamp := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	blockDependencies := blockchain.BlockDependencies{
		// use the custom clock function to get the same blocks
		Clock: func() time.Time {
			timestamp = timestamp.Add(time.Hour)
			return timestamp
		},
		Proofer: proofers.ProofOfWork{TargetBit: 248},
	}

	genesisBlock, err := blockchain.NewGenesisBlockEx(
		context.Background(),
		blockchain.NewGenesisBlockExParams{
			Clock:   blockDependencies.Clock,
			Data:    blockchain.NewData("genesis block"),
			Proofer: blockDependencies.Proofer,
		},
	)
	if err != nil {
		log.Fatalf("unable to create a new genesis block: %s", err)
	}

	blocks := []blockchain.Block{genesisBlock}
	for i := 0; i < 5; i++ {
		block, err := blockchain.NewBlockEx(
			context.Background(),
			blockchain.NewBlockExParams{
				Clock:     blockDependencies.Clock,
				Data:      blockchain.NewData(fmt.Sprintf("block #%d", i)),
				PrevBlock: blocks[len(blocks)-1],
				Proofer:   blockDependencies.Proofer,
			},
		)
		if err != nil {
			log.Fatalf("unable to create a new block: %s", err)
		}

		blocks = append(blocks, block)
	}

	blocksBytes, _ := json.MarshalIndent(blocks, "", "  ")
	fmt.Println(string(blocksBytes))

	// Output:
	// [
	//   {
	//     "Timestamp": "2006-01-02T16:04:05Z",
	//     "Data": "genesis block",
	//     "Hash": "248:225:00e26abd9974fcdea4b32eca43c9dc5c67fffa8efd53cebffa9b049fd6c2bb36",
	//     "PrevHash": ""
	//   },
	//   {
	//     "Timestamp": "2006-01-02T17:04:05Z",
	//     "Data": "block #0",
	//     "Hash": "248:198:0058f5dae6ca3451801a276c94862c7cce085e6f9371e50d80ddbb87c1438faf",
	//     "PrevHash": "248:225:00e26abd9974fcdea4b32eca43c9dc5c67fffa8efd53cebffa9b049fd6c2bb36"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T18:04:05Z",
	//     "Data": "block #1",
	//     "Hash": "248:15:002fc891ad012c4a89f7b267a2ec1767415c627ff69b88b90a93be938b026efa",
	//     "PrevHash": "248:198:0058f5dae6ca3451801a276c94862c7cce085e6f9371e50d80ddbb87c1438faf"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T19:04:05Z",
	//     "Data": "block #2",
	//     "Hash": "248:136:003c7def3d467a759fad481c03cadbd62e62b2c5dbc10e4bbb6e1944c158a8be",
	//     "PrevHash": "248:15:002fc891ad012c4a89f7b267a2ec1767415c627ff69b88b90a93be938b026efa"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T20:04:05Z",
	//     "Data": "block #3",
	//     "Hash": "248:65:00d5800e119abe44d89469c2161be7f9645d7237697c6d14b4a72717893582fa",
	//     "PrevHash": "248:136:003c7def3d467a759fad481c03cadbd62e62b2c5dbc10e4bbb6e1944c158a8be"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T21:04:05Z",
	//     "Data": "block #4",
	//     "Hash": "248:173:00b6863763acd6ec77ca3521589d8e68c118efe855657d702783e8e6aee169a9",
	//     "PrevHash": "248:65:00d5800e119abe44d89469c2161be7f9645d7237697c6d14b4a72717893582fa"
	//   }
	// ]
}

func ExampleBlockGroup() {
	timestamp := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	blockChunks := []blockchain.BlockGroup{
		// chunk #0
		{
			{
				Timestamp: timestamp.Add(6 * time.Hour),
				Data:      blockchain.NewData("block #4"),
				Hash: "248:" +
					"173:" +
					"00b6863763acd6ec77ca3521589d8e68c118efe855657d702783e8e6aee169a9",
				PrevHash: "248:" +
					"65:" +
					"00d5800e119abe44d89469c2161be7f9645d7237697c6d14b4a72717893582fa",
			},
			{
				Timestamp: timestamp.Add(5 * time.Hour),
				Data:      blockchain.NewData("block #3"),
				Hash: "248:" +
					"65:" +
					"00d5800e119abe44d89469c2161be7f9645d7237697c6d14b4a72717893582fa",
				PrevHash: "248:" +
					"136:" +
					"003c7def3d467a759fad481c03cadbd62e62b2c5dbc10e4bbb6e1944c158a8be",
			},
		},

		// chunk #1
		{
			{
				Timestamp: timestamp.Add(4 * time.Hour),
				Data:      blockchain.NewData("block #2"),
				Hash: "248:" +
					"136:" +
					"003c7def3d467a759fad481c03cadbd62e62b2c5dbc10e4bbb6e1944c158a8be",
				PrevHash: "248:" +
					"15:" +
					"002fc891ad012c4a89f7b267a2ec1767415c627ff69b88b90a93be938b026efa",
			},
			{
				Timestamp: timestamp.Add(3 * time.Hour),
				Data:      blockchain.NewData("block #1"),
				Hash: "248:" +
					"15:" +
					"002fc891ad012c4a89f7b267a2ec1767415c627ff69b88b90a93be938b026efa",
				PrevHash: "248:" +
					"198:" +
					"0058f5dae6ca3451801a276c94862c7cce085e6f9371e50d80ddbb87c1438faf",
			},
		},

		// chunk #2
		{
			{
				Timestamp: timestamp.Add(2 * time.Hour),
				Data:      blockchain.NewData("block #0"),
				Hash: "248:" +
					"198:" +
					"0058f5dae6ca3451801a276c94862c7cce085e6f9371e50d80ddbb87c1438faf",
				PrevHash: "248:" +
					"225:" +
					"00e26abd9974fcdea4b32eca43c9dc5c67fffa8efd53cebffa9b049fd6c2bb36",
			},
			{
				Timestamp: timestamp.Add(time.Hour),
				Data:      blockchain.NewData("genesis block"),
				Hash: "248:" +
					"225:" +
					"00e26abd9974fcdea4b32eca43c9dc5c67fffa8efd53cebffa9b049fd6c2bb36",
				PrevHash: "",
			},
		},
	}

	var prependedChunk blockchain.BlockGroup
	proofer := proofers.ProofOfWork{TargetBit: 248}
	for index, blockChunk := range blockChunks {
		validationMode := blockchain.AsBlockchainChunk
		if index == len(blockChunks)-1 {
			validationMode = blockchain.AsFullBlockchain
		}

		err := blockChunk.IsValid(prependedChunk, validationMode, proofer)
		if err != nil {
			log.Fatalf("chunk #%d is incorrect: %v", index, err)
		}

		prependedChunk = blockChunk
	}

	fmt.Println("all chunks are correct")

	// Output:
	// all chunks are correct
}
