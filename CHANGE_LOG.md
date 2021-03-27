# Change Log

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
