# go-blockchain

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-blockchain?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-blockchain)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-blockchain)](https://goreportcard.com/report/github.com/thewizardplusplus/go-blockchain)
[![lint](https://github.com/thewizardplusplus/go-blockchain/actions/workflows/lint.yaml/badge.svg)](https://github.com/thewizardplusplus/go-blockchain/actions/workflows/lint.yaml)
[![test](https://github.com/thewizardplusplus/go-blockchain/actions/workflows/test.yaml/badge.svg)](https://github.com/thewizardplusplus/go-blockchain/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-blockchain/graph/badge.svg?token=1SW1UEPF0I)](https://codecov.io/gh/thewizardplusplus/go-blockchain)

The library that implements models and algorithms of blockchain.

## Features

- models:
  - block data:
    - operations:
      - conversion to a string;
      - comparison for equality with another block data;
    - wrappers:
      - wrapper that adds support for the following operations to those block data that cannot do them:
        - conversion to a string:
          - implementation of the `fmt.Stringer` interface;
          - implementation of the `encoding.TextMarshaler` interface;
        - comparison for equality with another block data;
  - block:
    - storing:
      - timestamp;
      - block data;
      - hash;
      - previous hash;
    - operations:
      - creation (using a proofer);
      - getting merged data;
      - comparison for equality with another block;
      - self-validation (using a proofer);
  - genesis block:
    - based on a usual block without a previous hash;
  - block group:
    - storing:
      - group of blocks;
    - operations:
      - self-validation (using a proofer):
        - modes:
          - as a full blockchain;
          - as a blockchain chunk;
        - takes into account a prepended chunk;
        - allows empty block groups;
      - validation of the last block (using a proofer):
        - modes:
          - as a full blockchain;
          - as a blockchain chunk;
      - search of differences between two block groups:
        - returns lengths of different prefixes of the compared block groups;
        - based on a hash table index;
      - calculating a total difficulty of blocks;
  - block group loaders:
    - loading block groups via the external interface;
    - automatically saving the loaded block groups to a storage;
    - search of differences between two block group loaders:
      - loads and compares only one block chunk from every block group loader;
    - wrappers:
      - chunk validating loader:
        - automatically validates the loaded block group as a blockchain chunk;
      - last block validating loader:
        - automatically validates the last block from the loaded block group;
        - automatically preloads the next block group to perform the above validation;
      - memoizing loader:
        - remembers loaded block groups;
        - restricts the quantity of the remembered block groups:
          - stores the loaded block groups in the LRU cache;
    - kinds:
      - memory loader:
        - loading blocks from the block group;
  - blockchain:
    - storing:
      - storage;
      - last block;
    - operations:
      - creation:
        - loading the last block from the storage;
        - when the storage is empty (optional):
          - creation a genesis block using a proofer;
          - storing the genesis block to the storage;
      - loading block groups from the storage;
      - adding a block:
        - creation a block using a proofer;
        - storing the block to the storage;
      - merging with another blockchain:
        - selecting a fork based on a maximal total difficulty;
        - with automatic deleting orphan blocks;
- proofers:
  - operations:
    - block hashing;
    - block difficulty calculating;
    - block validation;
  - kinds:
    - [proof of work](https://en.wikipedia.org/wiki/Proof_of_work):
      - based on the [Hashcash](https://en.wikipedia.org/wiki/Hashcash) algorithm;
      - additional storing in a block (in a hash actually):
        - nonce;
        - target bit;
      - difficulty is defined as an inverse target bit;
- storages:
  - operations:
    - creation from a block group;
    - loading block groups;
    - loading the last block;
    - storing a block;
    - storing a block group (optional);
    - deleting a block;
    - deleting a block group (optional);
  - wrappers:
    - wrapper that adds support for the following operations to those storages that cannot do them:
      - storing a block group;
      - deleting a block group;
  - kinds:
    - memory storage:
      - storing blocks in memory.

## Installation

```
$ go get github.com/thewizardplusplus/go-blockchain
```

## Examples

`blockchain.Blockchain`:

```go
package main

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

func main() {
	timestamp := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	blockDependencies := blockchain.BlockDependencies{
		// use the custom clock function to get the same blocks
		Clock: func() time.Time {
			timestamp = timestamp.Add(time.Hour)
			return timestamp
		},
		Proofer: proofers.ProofOfWork{
			TargetBit: 248,
		},
	}

	blockchainInstance, err := blockchain.NewBlockchainEx(
		context.Background(),
		blockchain.NewBlockchainExParams{
			Dependencies: blockchain.Dependencies{
				BlockDependencies: blockDependencies,
				Storage:           storing.NewGroupStorage(&storages.MemoryStorage{}),
			},
			GenesisBlockData: mo.Some(blockchain.NewData("genesis block")),
		},
	)
	if err != nil {
		log.Fatalf("unable to create a new blockchain: %v", err)
	}

	const blockCount = 5
	for i := 0; i < blockCount; i++ {
		if err := blockchainInstance.AddBlockEx(
			context.Background(),
			blockchain.NewData(fmt.Sprintf("block #%d", i)),
		); err != nil {
			log.Fatalf("unable to add a new block: %v", err)
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
```

`blockchain.Blockchain.Merge()`:

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/samber/mo"
	blockchain "github.com/thewizardplusplus/go-blockchain"
	"github.com/thewizardplusplus/go-blockchain/proofers"
	"github.com/thewizardplusplus/go-blockchain/storing"
	"github.com/thewizardplusplus/go-blockchain/storing/storages"
)

func main() {
	blockDependencies := blockchain.BlockDependencies{
		Clock: time.Now,
		Proofer: proofers.ProofOfWork{
			TargetBit: 248,
		},
	}

	timestamp := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	blockGroupOne := blockchain.BlockGroup{
		{
			Timestamp: timestamp.Add(2*time.Hour + 40*time.Minute),
			Data:      blockchain.NewData("block #1.2"),
			Hash: "250:" +
				"57:" +
				"02988cf0f90c7771e726245c71416c6c376a2a473e90f3d0a7ba9787af421b34",
			PrevHash: "250:" +
				"7:" +
				"031d530789698389a084fd7a32e4b315d59fb0791a7b22ac4dce90be5a030eb5",
		},
		{
			Timestamp: timestamp.Add(2*time.Hour + 20*time.Minute),
			Data:      blockchain.NewData("block #1.1"),
			Hash: "250:" +
				"7:" +
				"031d530789698389a084fd7a32e4b315d59fb0791a7b22ac4dce90be5a030eb5",
			PrevHash: "240:" +
				"25578:" +
				"0000d382b7d47324d79ba6178449f9ebbd08a20412c2fa548a32f6ad217f6ce9",
		},
		{
			Timestamp: timestamp.Add(time.Hour),
			Data:      blockchain.NewData("block #0"),
			Hash: "240:" +
				"25578:" +
				"0000d382b7d47324d79ba6178449f9ebbd08a20412c2fa548a32f6ad217f6ce9",
			PrevHash: "240:" +
				"73021:" +
				"00004a15cf538f5e4d3592c68ee4ac6dd3d3b99d7fa5effcd75fe07c58eb213e",
		},
		{
			Timestamp: timestamp,
			Data:      blockchain.NewData("genesis block"),
			Hash: "240:" +
				"73021:" +
				"00004a15cf538f5e4d3592c68ee4ac6dd3d3b99d7fa5effcd75fe07c58eb213e",
			PrevHash: "",
		},
	}
	blockchainInstanceOne, err := blockchain.NewBlockchainEx(
		context.Background(),
		blockchain.NewBlockchainExParams{
			Dependencies: blockchain.Dependencies{
				BlockDependencies: blockDependencies,
				Storage: storing.NewGroupStorage(
					storages.NewMemoryStorage(blockGroupOne),
				),
			},
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
			Hash: "240:" +
				"885:" +
				"0000afa95e15291e5d6e7b5454292841114904f4d4b81c8187e838b7fe7d7b25",
			PrevHash: "240:" +
				"25578:" +
				"0000d382b7d47324d79ba6178449f9ebbd08a20412c2fa548a32f6ad217f6ce9",
		},
		{
			Timestamp: timestamp.Add(time.Hour),
			Data:      blockchain.NewData("block #0"),
			Hash: "240:" +
				"25578:" +
				"0000d382b7d47324d79ba6178449f9ebbd08a20412c2fa548a32f6ad217f6ce9",
			PrevHash: "240:" +
				"73021:" +
				"00004a15cf538f5e4d3592c68ee4ac6dd3d3b99d7fa5effcd75fe07c58eb213e",
		},
		{
			Timestamp: timestamp,
			Data:      blockchain.NewData("genesis block"),
			Hash: "240:" +
				"73021:" +
				"00004a15cf538f5e4d3592c68ee4ac6dd3d3b99d7fa5effcd75fe07c58eb213e",
			PrevHash: "",
		},
	}
	blockchainInstanceTwo, err := blockchain.NewBlockchainEx(
		context.Background(),
		blockchain.NewBlockchainExParams{
			Dependencies: blockchain.Dependencies{
				BlockDependencies: blockDependencies,
				Storage: storing.NewGroupStorage(
					storages.NewMemoryStorage(blockGroupTwo),
				),
			},
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
	//     "Hash": "240:885:0000afa95e15291e5d6e7b5454292841114904f4d4b81c8187e838b7fe7d7b25",
	//     "PrevHash": "240:25578:0000d382b7d47324d79ba6178449f9ebbd08a20412c2fa548a32f6ad217f6ce9"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T16:04:05Z",
	//     "Data": "block #0",
	//     "Hash": "240:25578:0000d382b7d47324d79ba6178449f9ebbd08a20412c2fa548a32f6ad217f6ce9",
	//     "PrevHash": "240:73021:00004a15cf538f5e4d3592c68ee4ac6dd3d3b99d7fa5effcd75fe07c58eb213e"
	//   },
	//   {
	//     "Timestamp": "2006-01-02T15:04:05Z",
	//     "Data": "genesis block",
	//     "Hash": "240:73021:00004a15cf538f5e4d3592c68ee4ac6dd3d3b99d7fa5effcd75fe07c58eb213e",
	//     "PrevHash": ""
	//   }
	// ]
}
```

`blockchain.Block`:

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/samber/mo"
	"github.com/thewizardplusplus/go-blockchain"
	"github.com/thewizardplusplus/go-blockchain/proofers"
)

func main() {
	timestamp := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	blockDependencies := blockchain.BlockDependencies{
		// use the custom clock function to get the same blocks
		Clock: func() time.Time {
			timestamp = timestamp.Add(time.Hour)
			return timestamp
		},
		Proofer: proofers.ProofOfWork{
			TargetBit: 248,
		},
	}

	genesisBlock, err := blockchain.NewGenesisBlockEx(
		context.Background(),
		blockchain.NewGenesisBlockExParams{
			Dependencies: blockDependencies,
			Data:         blockchain.NewData("genesis block"),
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
				Dependencies: blockDependencies,
				Data:         blockchain.NewData(fmt.Sprintf("block #%d", i)),
				PrevBlock:    mo.Some(blocks[len(blocks)-1]),
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
```

`blockchain.BlockGroup`:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/thewizardplusplus/go-blockchain"
	"github.com/thewizardplusplus/go-blockchain/proofers"
)

func main() {
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
```

`loading.LoadStorage()`:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/thewizardplusplus/go-blockchain"
	"github.com/thewizardplusplus/go-blockchain/loading"
	"github.com/thewizardplusplus/go-blockchain/loading/loaders"
	"github.com/thewizardplusplus/go-blockchain/proofers"
	"github.com/thewizardplusplus/go-blockchain/storing"
	"github.com/thewizardplusplus/go-blockchain/storing/storages"
)

type LoggingLoader struct {
	Loader blockchain.Loader
}

func (loader LoggingLoader) LoadBlocks(cursor interface{}, count int) (
	blocks blockchain.BlockGroup,
	nextCursor interface{},
	err error,
) {
	fmt.Printf("[DEBUG] load the blocks corresponding to cursor %v\n", cursor)

	return loader.Loader.LoadBlocks(cursor, count)
}

func main() {
	timestamp := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	blocks := blockchain.BlockGroup{
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
	}

	var storage storages.MemoryStorage
	proofer := proofers.ProofOfWork{TargetBit: 248}
	if _, err := loading.LoadStorage(
		storing.NewGroupStorage(&storage),
		loading.LastBlockValidatingLoader{
			Loader: loading.NewMemoizingLoader(1, loading.ChunkValidatingLoader{
				Loader: LoggingLoader{
					Loader: loaders.MemoryLoader(blocks),
				},
				Proofer: proofer,
			}),
			Proofer: proofer,
		},
		nil,
		2,
	); err != nil {
		log.Fatalf("unable to load the blocks: %v", err)
	}

	loadedBlocks, _, _ := storage.LoadBlocks(nil, len(blocks))
	blocksBytes, _ := json.MarshalIndent(loadedBlocks, "", "  ")
	fmt.Println(string(blocksBytes))

	// Output:
	// [DEBUG] load the blocks corresponding to cursor <nil>
	// [DEBUG] load the blocks corresponding to cursor 2
	// [DEBUG] load the blocks corresponding to cursor 4
	// [DEBUG] load the blocks corresponding to cursor 6
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
```

## License

The MIT License (MIT)

Copyright &copy; 2021-2022, 2025 thewizardplusplus
