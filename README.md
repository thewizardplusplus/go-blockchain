# go-blockchain

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-blockchain?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-blockchain)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-blockchain)](https://goreportcard.com/report/github.com/thewizardplusplus/go-blockchain)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-blockchain.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-blockchain)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-blockchain/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-blockchain)

The library that implements models and algorithms of blockchain.

## Features

- models:
  - block:
    - storing:
      - timestamp;
      - custom data;
      - hash;
      - previous hash;
    - operations:
      - creation (using a proofer);
      - getting merged data;
      - self-validation (using a proofer);
  - genesis block:
    - based on a usual block without a previous hash;
- proofers:
  - operations:
    - block hashing;
    - block validation;
  - kinds:
    - simple:
      - based on once hashing by the [SHA-256](https://en.wikipedia.org/wiki/SHA-2) algorithm;
    - [proof of work](https://en.wikipedia.org/wiki/Proof_of_work):
      - based on the [Hashcash](https://en.wikipedia.org/wiki/Hashcash) algorithm;
      - additional storing in a block (in a hash actually):
        - nonce;
        - target bit.

## Installation

Prepare the directory:

```
$ mkdir --parents "$(go env GOPATH)/src/github.com/thewizardplusplus/"
$ cd "$(go env GOPATH)/src/github.com/thewizardplusplus/"
```

Clone this repository:

```
$ git clone https://github.com/thewizardplusplus/go-blockchain.git
$ cd go-blockchain
```

Install dependencies with the [dep](https://golang.github.io/dep/) tool:

```
$ dep ensure -vendor-only
```

## Examples

`blockchain.Blockchain`:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/thewizardplusplus/go-blockchain"
	"github.com/thewizardplusplus/go-blockchain/proofers"
	"github.com/thewizardplusplus/go-blockchain/storages"
)

type StringHasher string

func (hasher StringHasher) Hash() string {
	return string(hasher)
}

func main() {
	timestamp := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	blockDependencies := blockchain.BlockDependencies{
		// use the custom clock function to get the same blocks
		Clock: func() time.Time {
			timestamp = timestamp.Add(time.Hour)
			return timestamp
		},
		Proofer: proofers.ProofOfWork{TargetBit: 248},
	}

	var storage storages.MemoryStorage
	blockchain, err := blockchain.NewBlockchain(
		StringHasher("genesis block"),
		blockchain.Dependencies{
			BlockDependencies: blockDependencies,
			Storage:           &storage,
		},
	)
	if err != nil {
		log.Fatalf("unable to create the blockchain: %v", err)
	}

	for i := 0; i < 5; i++ {
		if err := blockchain.AddBlock(
			StringHasher(fmt.Sprintf("block #%d", i)),
		); err != nil {
			log.Fatalf("unable to add the block: %v", err)
		}
	}

	blocksBytes, _ := json.MarshalIndent(storage.Blocks(), "", "  ")
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

`blockchain.Block`:

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/thewizardplusplus/go-blockchain"
	"github.com/thewizardplusplus/go-blockchain/proofers"
)

type StringHasher string

func (hasher StringHasher) Hash() string {
	return string(hasher)
}

func main() {
	timestamp := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	dependencies := blockchain.Dependencies{
		// use the custom clock function to get the same blocks
		Clock: func() time.Time {
			timestamp = timestamp.Add(time.Hour)
			return timestamp
		},
		Proofer: proofers.ProofOfWork{TargetBit: 248},
	}

	blocks := []blockchain.Block{
		blockchain.NewGenesisBlock(StringHasher("genesis block"), dependencies),
	}
	for i := 0; i < 5; i++ {
		blocks = append(blocks, blockchain.NewBlock(
			StringHasher(fmt.Sprintf("block #%d", i)),
			blocks[len(blocks)-1],
			dependencies,
		))
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

## License

The MIT License (MIT)

Copyright &copy; 2021 thewizardplusplus
