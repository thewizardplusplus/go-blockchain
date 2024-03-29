# Change Log

## [v1.4](https://github.com/thewizardplusplus/go-blockchain/tree/v1.4) (2022-03-09)

Add merging the blockchain model with another one.

- models:
  - blockchain:
    - operations:
      - merging with another blockchain:
        - selecting a fork based on a maximal total difficulty;
        - with automatic deleting orphan blocks;
- storages:
  - operations:
    - creation from a block group;
- additionally:
  - return nil instead empty array from `loaders.MemoryLoader.LoadBlocks()` method;
  - copy the loaded blocks in the `storages.MemoryStorage.LoadBlocks()` method;
  - fix the unit tests:
    - add the missed tests without blocks;
    - add the tests with checking pointers for the `storages.MemoryStorage.LoadBlocks()` method;
  - examples:
    - add the example with merging the blockchain model with another one.

## [v1.4-beta.1](https://github.com/thewizardplusplus/go-blockchain/tree/v1.4-beta.1) (2022-03-06)

Add calculating a difficulty of blocks.

- models:
  - block group:
    - operations:
      - calculating a difficulty of blocks;
- proofers:
  - operations:
    - block difficulty calculating;
  - kinds:
    - [proof of work](https://en.wikipedia.org/wiki/Proof_of_work):
      - difficulty is defined as an inverse target bit.

## [v1.4-beta](https://github.com/thewizardplusplus/go-blockchain/tree/v1.4-beta) (2022-03-05)

Add a search of differences between two block groups.

- models:
  - block group:
    - operations:
      - search of differences between two block groups:
        - returns lengths of different prefixes of the compared block groups;
        - based on a hash table index;
  - block group loaders:
    - search of differences between two block group loaders:
      - loads and compares only one block chunk from every block group loader;
- remove the `blockchain.BlockGroup.FindBlock()` method.

## [v1.4-alpha.2](https://github.com/thewizardplusplus/go-blockchain/tree/v1.4-alpha.2) (2022-02-11)

Add a search of a block index and fix the bug with sorting in the `storages.MemoryStorage` structure.

- models:
  - block group:
    - operations:
      - search of a block index:
        - based on the [binary search](https://en.wikipedia.org/wiki/Binary_search_algorithm) algorithm;
- fix the bug with sorting in the `storages.MemoryStorage` structure:
  - fix the `storages.MemoryStorage.sortIfNeed()` method;
  - fix the `storages.MemoryStorage.DeleteBlock()` method;
  - fix the unit tests:
    - of the `storages.MemoryStorage.LoadBlocks()` method;
    - of the `storages.MemoryStorage.LoadLastBlock()` method;
    - of the `storages.MemoryStorage.StoreBlock()` method;
    - of the `storages.MemoryStorage.DeleteBlock()` method.

## [v1.4-alpha.1](https://github.com/thewizardplusplus/go-blockchain/tree/v1.4-alpha.1) (2021-12-26)

Add the model of block data.

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
        - comparison for equality with another block data.

## [v1.4-alpha](https://github.com/thewizardplusplus/go-blockchain/tree/v1.4-alpha) (2021-12-21)

Add to the block model the comparison for equality with another block; add to the storage the deleting of a block.

- models:
  - block:
    - operations:
      - comparison for equality with another block;
- storages:
  - operations:
    - deleting a block;
    - deleting a block group (optional);
  - wrappers:
    - wrapper that adds support for deleting a block group to those storages that cannot do this.

## [v1.3.2](https://github.com/thewizardplusplus/go-blockchain/tree/v1.3.2) (2021-11-24)

Add the loading of block groups to the storages (including the memory storage) and the blockchain model; simplify the passing of the proofer to the validators; make optional the creation of a genesis block on the blockchain creation; remove the redundant features.

- models:
  - blockchain:
    - operations:
      - creation:
        - make optional:
          - when the storage is empty:
            - creation a genesis block using a proofer;
            - storing the genesis block to the storage;
      - loading block groups from the storage;
- proofers:
  - kinds:
    - remove the `proofers.Simple` structure;
- storages:
  - operations:
    - loading block groups;
  - kinds:
    - memory storage:
      - remove the `storages.MemoryStorage.Blocks()` method;
- additionally:
  - replace the `blockchain.BlockDependencies` structure to the `blockchain.Proofer` interface:
    - in the `blockchain.Block` structure:
      - in the `blockchain.Block.IsValid()` method;
      - in the `blockchain.Block.IsValidGenesisBlock()` method;
    - in the `blockchain.BlockGroup` structure:
      - in the `blockchain.BlockGroup.IsValid()` method;
      - in the `blockchain.BlockGroup.IsLastBlockValid()` method;
    - in the `loading.ChunkValidatingLoader` structure;
    - in the `loading.LastBlockValidatingLoader` structure.

## [v1.3.1](https://github.com/thewizardplusplus/go-blockchain/tree/v1.3.1) (2021-10-24)

Add the memory loader; restrict the cache size in the memoizing loader; optimize the memory storage; replace the `blockchain.Hasher` interface to the `fmt.Stringer` interface; return the errors from the proofers; split the code to the additional packages.

- models:
  - block group loaders:
    - wrappers:
      - memoizing loader:
        - restricts the quantity of the remembered block groups:
          - stores the loaded block groups in the LRU cache;
    - kinds:
      - memory loader:
        - loading blocks from the block group;
- additionally:
  - replace the `blockchain.Hasher` interface to the `fmt.Stringer` interface;
  - return the error instead of a boolean flag from the `blockchain.Proofer.Validate()` method;
  - add the packages:
    - `storing`;
    - `loading`;
  - optimize the `storages.MemoryStorage` structure:
    - revert the use of the `container/heap` package;
    - optimize the detection of the last block;
    - optimize the sorting of the stored blocks;
  - examples:
    - fix the examples in the `README.md` file;
    - improve the example with the loading of the blocks.

## [v1.3](https://github.com/thewizardplusplus/go-blockchain/tree/v1.3) (2021-09-21)

Implementing the loading of block groups to a storage via the external interface; supporting the automatical validation of the loaded block groups.

- models:
  - block group:
    - operations:
      - validation of the last block (using a proofer):
        - modes:
          - as a full blockchain;
          - as a blockchain chunk;
  - block group loaders:
    - loading block groups via the external interface;
    - automatically saving the loaded block groups to a storage;
    - wrappers:
      - chunk validating loader:
        - automatically validates the loaded block group as a blockchain chunk;
      - last block validating loader:
        - automatically validates the last block from the loaded block group;
        - automatically preloads the next block group to perform the above validation;
      - memoizing loader:
        - remembers loaded block groups;
- storages:
  - operations:
    - storing a block group (optional);
  - wrappers:
    - wrapper that adds support for saving a block group to those storages that cannot do this.

## [v1.2](https://github.com/thewizardplusplus/go-blockchain/tree/v1.2) (2021-05-09)

Implementing a block group abstraction with self-validation; supporting the two modes of its validation: as a full blockchain and as a blockchain chunk.

- models:
  - block group:
    - storing:
      - group of blocks;
    - operations:
      - self-validation (using a proofer):
        - modes:
          - as a full blockchain;
          - as a blockchain chunk;
        - takes into account a prepended chunk;
        - allows empty block groups.

## [v1.1](https://github.com/thewizardplusplus/go-blockchain/tree/v1.1) (2021-03-27)

Implementing a blockchain abstraction that works with a storage and a memory storage for it.

- models:
  - blockchain:
    - storing:
      - storage;
      - last block;
    - operations:
      - creation:
        - loading the last block from the storage;
        - when the storage is empty:
          - creation a genesis block using a proofer;
          - storing the genesis block to the storage;
      - adding a block:
        - creation a block using a proofer;
        - storing the block to the storage;
- storages:
  - operations:
    - loading the last block;
    - storing a block;
  - kinds:
    - memory storage:
      - storing blocks in memory;
      - additional operations:
        - getting the stored blocks.

## [v1.0](https://github.com/thewizardplusplus/go-blockchain/tree/v1.0) (2021-03-13)

Major version. Implementing blocks and their validators (including the [proof of work](https://en.wikipedia.org/wiki/Proof_of_work) algorithm).
