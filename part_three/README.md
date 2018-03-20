# Part Three

[building-blockchain-in-go Part 3 Persistence and CLI](https://jeiwan.cc/posts/building-blockchain-in-go-part-3/)

本节将解决数据持久化（数据库相关）的问题，以及构建一个命令行工具来操控区块链。
区块链的数据库是分布式的，但是本节将暂时不考虑"分布式"，只关心存储问题。

比特币用的是 LevelDB，我们将采用 BoltDB：
1. 轻量级
2. 是用 Go 实现的
3. 不需要一个服务器就能运行
4. 允许我们构建我们想要的数据结构
>Bolt is a pure Go key/value store inspired by Howard Chu’s LMDB project.
The goal of the project is to provide a simple, fast, and reliable database
for projects that don’t require a full database server such as Postgres or MySQL.

BoltDB是没有数据类型的，它的keys和values都是byte array。
然后我们采用了 encoding/gob 来序列化我们的 Block 。

当然也可以使用 JSON，XML，Protocol Buffers等。

# Data structure
实施数据持久化之前，先想好如何把数据存到数据库里。这里参考比特币的做法。

Bitcoin Core uses two "buckets" to store data:
1. **blocks** stores metadata describing all the blocks in a chain
2. **chainstate** stores the state of a chain,
which is all currently unspent transaction outputs and some metadata.

当然，区块也会作为独立文件存储于磁盘上，这是基于性能上的考量：
读取某一个区块不需要加载所有的区块到内存里。

**我们目前不实现这一特性（也就是说，所有数据库都存在于一个文件里）**

在 **blocks** 里，键值对包括：
1. 'b' + 32-byte block hash -> block index record
2. 'f' + 4-byte file number -> file information record
3. 'l' -> 4-byte file number: the last block file number used
4. 'R' -> 1-byte boolean: whether we're in the process of reindexing
5. 'F' + 1-byte flag name length + flag name string -> 1 byte boolean: various flags that can be on or off
6. 't' + 32-byte transaction hash -> transaction index record

在 **chainstate**里，键值对包括：
1. 'c' + 32-byte transaction hash -> unspent transaction output record for that transaction
2. 'B' -> 32-byte block hash: the block hash up to which the database represents the unspent transaction outputs

因为没有transactions，所以只需要 blocks bucket。

因为整个数据库就是一个文件，所以不需要file number 相关的数据。

这样的话，我们只需要：
1. 32-byte block-hash -> Block structure (serialized)
2. 'l' -> the hash of the last block in a chain

# Serialization（序列化）
BoltDB的 value 只能是 []byte

# Persistence（持久化）
也就是存储于磁盘上。需要我们对 NewBlockchain 进行改造。
1. Open a DB file
2. Check if there's a blockchain stored in it
3. if there's a blockchain:
    - Create a new **Blockchain** instance
    - Set the tip of the **Blockchain** instance to the last block hash stored in the DB
4. if there's no existing blockchain:
    - Create the genesis block
    - Store in the DB
    - Save the genesis block's hash as the last block hash
    - Create a new Blockchain instance with its tip pointing at the genesis block

# How to run the code
```
➜  part_three git:(master) ✗ go run *.go printchain
Prev. hash:
Data: Genesis Block
Hash: 000000cb05cc9c46dfe4927497e2818e3406992c211a124b02f661b12279d1a2
PoW: true
➜  part_three git:(master) ✗ go run *.go addblock -data "Send 1 BTC to Ivan"
Mining the block containing "Send 1 BTC to Ivan"
000000295b04ae065a761b2dec943f1cc297cf261eae58ba317badd960b1e559

Success!
➜  part_three git:(master) ✗ go run *.go addblock -data "Pay 0.31337 BTC for a coffee"
Mining the block containing "Pay 0.31337 BTC for a coffee"
000000fdfdb460d52223332160e94f6cb2b938e87dcb8f2ed815099bf5b64d47

Success!
➜  part_three git:(master) ✗ go run *.go printchain
Prev. hash: 000000295b04ae065a761b2dec943f1cc297cf261eae58ba317badd960b1e559
Data: Pay 0.31337 BTC for a coffee
Hash: 000000fdfdb460d52223332160e94f6cb2b938e87dcb8f2ed815099bf5b64d47
PoW: true

Prev. hash: 000000cb05cc9c46dfe4927497e2818e3406992c211a124b02f661b12279d1a2
Data: Send 1 BTC to Ivan
Hash: 000000295b04ae065a761b2dec943f1cc297cf261eae58ba317badd960b1e559
PoW: true

Prev. hash:
Data: Genesis Block
Hash: 000000cb05cc9c46dfe4927497e2818e3406992c211a124b02f661b12279d1a2
PoW: true

```

# Next
addresses, wallets, transactions