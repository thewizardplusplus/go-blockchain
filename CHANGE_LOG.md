# Change Log

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
