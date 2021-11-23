package loading_test

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

type StringData string

func (data StringData) String() string {
	return string(data)
}

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

func ExampleLoadStorage() {
	timestamp := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	blocks := blockchain.BlockGroup{
		{
			Timestamp: timestamp.Add(6 * time.Hour),
			Data:      StringData("block #4"),
			Hash: "248:" +
				"173:" +
				"00b6863763acd6ec77ca3521589d8e68c118efe855657d702783e8e6aee169a9",
			PrevHash: "248:" +
				"65:" +
				"00d5800e119abe44d89469c2161be7f9645d7237697c6d14b4a72717893582fa",
		},
		{
			Timestamp: timestamp.Add(5 * time.Hour),
			Data:      StringData("block #3"),
			Hash: "248:" +
				"65:" +
				"00d5800e119abe44d89469c2161be7f9645d7237697c6d14b4a72717893582fa",
			PrevHash: "248:" +
				"136:" +
				"003c7def3d467a759fad481c03cadbd62e62b2c5dbc10e4bbb6e1944c158a8be",
		},
		{
			Timestamp: timestamp.Add(4 * time.Hour),
			Data:      StringData("block #2"),
			Hash: "248:" +
				"136:" +
				"003c7def3d467a759fad481c03cadbd62e62b2c5dbc10e4bbb6e1944c158a8be",
			PrevHash: "248:" +
				"15:" +
				"002fc891ad012c4a89f7b267a2ec1767415c627ff69b88b90a93be938b026efa",
		},
		{
			Timestamp: timestamp.Add(3 * time.Hour),
			Data:      StringData("block #1"),
			Hash: "248:" +
				"15:" +
				"002fc891ad012c4a89f7b267a2ec1767415c627ff69b88b90a93be938b026efa",
			PrevHash: "248:" +
				"198:" +
				"0058f5dae6ca3451801a276c94862c7cce085e6f9371e50d80ddbb87c1438faf",
		},
		{
			Timestamp: timestamp.Add(2 * time.Hour),
			Data:      StringData("block #0"),
			Hash: "248:" +
				"198:" +
				"0058f5dae6ca3451801a276c94862c7cce085e6f9371e50d80ddbb87c1438faf",
			PrevHash: "248:" +
				"225:" +
				"00e26abd9974fcdea4b32eca43c9dc5c67fffa8efd53cebffa9b049fd6c2bb36",
		},
		{
			Timestamp: timestamp.Add(time.Hour),
			Data:      StringData("genesis block"),
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
