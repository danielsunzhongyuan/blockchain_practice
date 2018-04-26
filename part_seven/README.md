# Part Seven
[building-blockchain-in-go Part 7 Network](https://jeiwan.cc/posts/building-blockchain-in-go-part-7/)

目前为止，我们已经实现了区块链的所有核心功能：匿名、安全、随机产生的地址；
区块链数据存储；Proof-of-work系统；可靠的存储交易数据的方式。

这些特性都很关键，但是还少了一样 - **网络**！是网络使得数字加密货币如此的成功。
只有一个节点的话，区块链是没有意义的；只有一个用户的话，加密是没有意义的。

前面介绍的区块链的特性，可以当作一种规则，就像人与人之间相处时所遵循的规则一样。

**同时**：如果没有大量的网络节点，这些规则也时没用的！

>本系列的作者没时间实现一个真正的P2P的网络的。他只会展示一个最常见的场景，这个场景涉及到多种类型的节点。
除了这个场景的其他场景，就不保证能工作了。
所以，基于这个基本的场景，进而实现一个P2P网络，就当作一个挑战吧。


# 区块链网络
区块链的网络是去中心化的，也就是说，并没有一个中央服务器，也不存在客户端。
只有节点，而且每一个节点都是完整的。
一个节点是完整的，意思是：它既是服务器也是客户端。这一点和普通的Web应用不同。

区块链的网络是一个P2P（Pear-to-Pear）的网络，所有节点都是相互直连。没有上下级关系，完全扁平的拓扑结构。

![p2p-network](./p2p-network.png)

这就对每个节点的要求都很高，需要能够和其他节点交互、需要去获取其他节点的状态、和自己的状态做比较、更新自己的状态等等。

# 节点的职责
区块链的节点可以在网络里扮演不同的职责，比如：
1. 矿机（Miner）
纯挖矿，就为了解决PoW问题。但是比特币只是区块链的一种实现方式，其他的实现方式不一定就有PoW问题，也就不一定需要挖矿了，例如 Proof-of-Stake blockchain
2. Full node.
这些节点要校验区块、校验交易。它会包含所有的区块链的数据，也包含了所有节点的信息，可以帮助其他节点相互发现，因此有路由的作用。
网络里必须要有大量的这种 Full Nodes，正式这些node会决定某个区块、某个交易是不是合法的。
3. SPV - Simplified Payment Verification
这个前面已经提到了。这些节点不会存储区块链的完整副本，但是缺可以拿来校验交易。
一个SPV node往往依附于一个Full node，而一个Full node可以连接很多SPV node。
SPV node的作用是，极大地提升了校验的速度，从而使得钱包相关的概念得以实现和使用。
（没有SPV，查余额等操作就太慢了，那这个系统就没法用了）

# Network Simplification
P2P网络太复杂了，为了达到演示的目的，我们需要做很多简化的工作。
网络里需要很多节点，用虚拟机或者docker可以进行实验，但是也还是复杂，所以直接用端口进行区分即可。
例如：127.0.0.1:3000, 127.0.0.1:3001, 127.0.0.1:3002 就可以算作是3个节点了。

同时，对wallet、blockchain（数据库）文件也要做类似区分，也就是：blockchain_3000.db，blockchain_3001.db，wallet_3000.dat，wallet_3001.dat 。。。

# 实现
动手之前，先考虑一个问题：
当我们第一次下载所谓的"Bitcoin Core"并运行时，会发生什么？
我们需要连接到某些个节点上，然后下载最新的区块链数据。那么从哪些节点下载呢？

直接连接某一个具体的节点是不合适的，因为它可能会被攻击、会宕机。
因此，Bitcoin Core采用了 DNS seeds 的方式。它们不是节点，只是DNS服务器，指向了一些节点。

第一次下载Bitcoin Core时，先连接到一个 DNS seed，获取一些 Full node的地址，然后从这些Full node上获取数据。

在我们的实现里，我们暂时还是要中心化的，我们需要3个节点：
1. 中央节点：其他节点都会连接到这个节点上。
2. 矿机：这个节点会存储新的交易信息，当拥有足够的交易后，就会挖一个新的区块
3. A wallet node：这个节点会用作发送比特币。但是不像SPV节点，它会存储所有区块链数据。

# 最常见的场景
本章节要实现的场景是：
1. 中央节点产生区块链
2. Wallet节点连接到中央节点，然后下载区块链
3. 矿机节点连接到中央节点，然后下载区块链
4. Wallet节点产生交易（transaction）
5. 矿机节点接收交易，并放在内存里
6. 当内存里的交易数达到一定数量时，矿机会挖一个新的区块
7. 每新挖一个区块，都将发送给中央节点
8. Wallet节点则从中央节点进行同步
9. Wallet节点里的用户去检查交易是否成功

上面的流程就是比特币里最重要的使用场景了（use case）

## version
节点通过传递消息的方式进行沟通。当一个新节点开始运行是，它会从DNS seed里获取多个节点，
然后传给它们 `version` 信息：
```
type version struct {
    Version     int
    BestHeight  int
    AddrFrom    string
}
```

因为我们只有一个区块链的版本，所以，`Version`里并不存储重要的信息。
`BestHeight`放的是节点区块链的长度。`AddrFrom`放的是发送方的地址。

节点接收到 `version` 信息后能做什么呢？它会返回它自己的`version`信息。
这样做有点类似一次握手：交换其他信息前先进行一次握手。
不只是握手，`version`还可以用来查找一条更长的区块链。当节点收到version后，它会检查它自己的区块链是否比 BestHeight 长。
如果不的话，就会下载丢失的区块。

为了接收信息，我们需要一个服务器：
```
var nodeAddress string
var knownNodes = []string{"localhost:3000"}

func StartServer(nodeID, minerAddress string) {
    nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
    miningAddress = minerAddress
    ln, err := net.Listen(protocol, nodeAddress)
    defer ln.Close()

    bc := NewBlockchain(nodeID)

    if nodeAddress != knownNodes[0] {
        sendVersion(knownNodes[0], bc)
    }

    for {
        conn, err := ln.Accept()
        go handleConnection(conn, bc)
    }
}
```
首先，我们固定了中央节点的地址：每个节点都要连接到它。
`minerAddress`参数是指接收挖矿奖励的地址。
如果当前节点不是中央节点，那么需要发送version信息给中央节点，并检查自己的区块链是否过时了。

发送version信息的方法是 `sendVersion`，发送的就是当前节点的地址、版本、bestHeight。
```
func sendVersion(addr string, bc *Blockchain) {
    bestHeight := bc.GetBestHeight()
    payload := gobEncode(version{nodeVersion, bestHeight, nodeAddress})
    request := append(commandToBytes("version"), payload...)
    sendData(addr, request)
}
```
`commandToBytes`和`bytesToCommand`是一对相反的方法。command是不超过12个字节的字符串。
接收到消息后，把command抽取出来，再根据command进行下一步操作：
```
func handleConnection(conn net.Conn, bc *Blockchain) {
    request, err := ioutil.ReadAll(conn)
    command := bytesToCommand(request[:commandLength])
    fmt.Printf("Received %s command\n", command)

    switch command {
    ...
    case "version":
        handleVersion(request, bc)
    default:
        fmt.Println("Unknown command!")
    }

    conn.Close()
}

func handleVersion(request []byte, bc *Blockchain) {
    var buff bytes.Buffer
    var payload verzion

    buff.Write(request[commandLength:])
    dec := gob.NewDecoder(&buff)
    err := dec.Decode(&payload)

    myBestHeight := bc.GetBestHeight()
    foreignerBestHeight := payload.BestHeight

    if myBestHeight < foreignerBestHeight {
        sendGetBlocks(payload.AddrFrom)
    } else if myBestHeight > foreignerBestHeight {
        sendVersion(payload.AddrFrom, bc)
    }

    if !nodeIsKnown(payload.AddrFrom) {
        knownNodes = append(knownNodes, payload.AddrFrom)
    }
}
```
这里 bestHeight 的逻辑是：如果本节点的 bestHeight小于外面的bestHeight，就去外面更新区块；
否则的话，把自己的version信息发给中央节点（目的就是让别的节点去更新区块）。

如何获取/更新区块呢？
```
type getblocks struct {
    AddrFrom string
}

func handleGetBlocks(request []byte, bc *Blockchain) {
    ...
    blocks := bc.GetBlockHashes()
    sendInv(payload.AddrFrom, "block", blocks)
}
```
在我们的实现里，handleGetBlocks会获取全部的区块哈希值（all block hashes）

inv是用来告诉其他节点，当前节点有哪些区块和交易。当然，也不是全量信息，只是哈希值。
```
type inv struct {
    AddrFrom    string
    Type        string // block or transaction
    Items       [][]byte
}

func handleInv(request []byte, bc *Blockchain) {
    ...
    fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

    if payload.Type == "block" {
        blocksInTransit = payload.Items

        blockHash := payload.Items[0]
        sendGetData(payload.AddrFrom, "block", blockHash)

        newInTransit := [][]byte{}
        for _, b := range blocksInTransit {
            if bytes.Compare(b, blockHash) != 0 {
                newInTransit = append(newInTransit, b)
            }
        }
        blocksInTransit = newInTransit
    }

    if payload.Type == "tx" {
        txID := payload.Items[0]

        if mempool[hex.EncodeToString(txID)].ID == nil {
            sendGetData(payload.AddrFrom, "tx", txID)
        }
    }
}
```
handleInv要难的多：
接收到blocks hashes后，要保存到 blocksInTransit 变量里，从而记录有哪些已经下载了的block；
handleInv允许我们从不同的节点下载block（当然我们的实现不这么做）；
把区块放到 transit 状态后，我们发送 getdata 命令给"发送inv消息的节点"，并更新 blocksInTransit。
在真实的P2P网络中，我们是期望从不同的节点获取blocks的。

在我们的实现里，我们并不会发送带有多条哈希值的 inv 消息，而只是一条（第一条）。
```
if payload.Type == "tx" {
        txID := payload.Items[0]
```
如果这个tx不在 mempool 里，就 getdata。


`getdata`是为了获取特定的 block / transaction。
```
type getdata struct {
    AddrFrom    string
    Type        string
    ID          []byte
}

func handleGetData(request []byte, bc *Blockchain) {
    ...
    if payload.Type == "block" {
        block, err := bc.GetBlock([]byte(payload.ID))

        sendBlock(payload.AddrFrom, &block)
    }

    if payload.Type == "tx" {
        txID := hex.EncodeToString(payload.ID)
        tx := mempool[txID]

        sendTx(payload.AddrFrom, &tx)
    }
}
```
这个处理器还是很直白的：请求block，就返回block；请求transaction，就返回transaction。

上面讲了很多消息的传送、交互等等，真正存储数据的结构是：
```
type block struct {
    AddrFrom string
    Block    []byte
}

type tx struct {
    AddFrom     string
    Transaction []byte
}
```
处理 block 消息比较简单：
收到一个新的block，就放到我们的区块链里。如果需要下载更多区块，就从下载前一个区块的节点下载它们。
下载完所有需要的block后，Reindex UTXO set。
> TODO: 接收新区块的时候要进行校验
> TODO: 其实不应该reindex UTXO set，而是应该Update，因为实际的区块链太大了。

代码如下：
```
func handleBlock(request []byte, bc *Blockchain) {
    ...

    blockData := payload.Block
    block := DeserializeBlock(blockData)

    fmt.Println("Recevied a new block!")
    bc.AddBlock(block)

    fmt.Printf("Added block %x\n", block.Hash)

    if len(blocksInTransit) > 0 {
        blockHash := blocksInTransit[0]
        sendGetData(payload.AddrFrom, "block", blockHash)

        blocksInTransit = blocksInTransit[1:]
    } else {
        UTXOSet := UTXOSet{bc}
        UTXOSet.Reindex()
    }
}
```

处理transaction就非常麻烦了：
首先把新的交易放到内存池里，当然也是应该先校验的；
然后检查是不是中央节点，如果是的话就发送给其他节点；
如果是矿机节点：
    首先判断要至少有两个交易，且挖矿的地址数大于0（也就是矿机），才开始挖矿；
    然后校验所有交易；
    把校验后的交易，连同一个奖励的交易，放到区块链里，然后UTXOSet.Reindex()；
    然后，把这些交易从内存池里删除；
    再发给网络中的其他节点，发送 inv 消息。
代码如下：
```
func handleTx(request []byte, bc *Blockchain) {
    ...
    txData := payload.Transaction
    tx := DeserializeTransaction(txData)
    mempool[hex.EncodeToString(tx.ID)] = tx

    if nodeAddress == knownNodes[0] {
        for _, node := range knownNodes {
            if node != nodeAddress && node != payload.AddFrom {
                sendInv(node, "tx", [][]byte{tx.ID})
            }
        }
    } else {
        if len(mempool) >= 2 && len(miningAddress) > 0 {
        MineTransactions:
            var txs []*Transaction

            for id := range mempool {
                tx := mempool[id]
                if bc.VerifyTransaction(&tx) {
                    txs = append(txs, &tx)
                }
            }

            if len(txs) == 0 {
                fmt.Println("All transactions are invalid! Waiting for new ones...")
                return
            }

            cbTx := NewCoinbaseTX(miningAddress, "")
            txs = append(txs, cbTx)

            newBlock := bc.MineBlock(txs)
            UTXOSet := UTXOSet{bc}
            UTXOSet.Reindex()

            fmt.Println("New block is mined!")

            for _, tx := range txs {
                txID := hex.EncodeToString(tx.ID)
                delete(mempool, txID)
            }

            for _, node := range knownNodes {
                if node != nodeAddress {
                    sendInv(node, "block", [][]byte{newBlock.Hash})
                }
            }

            if len(mempool) > 0 {
                goto MineTransactions
            }
        }
    }
}
```

# 结果
我们再来梳理一遍最开始说定义的场景：
1. 打开一个终端，set NODE_ID = 3000 （export NODE_ID=3000）。
新建一个钱包，新建一个区块链
这个是中央节点。
- 复制blockchain，这里为 4.1 做准备 `cp blockchain_3000.db blockchain_genesis.db`
2. 打开一个新终端，set NODE_ID = 3001
然后新建多个钱包地址，例如 WALLET_1，WALLET_2，WALLET_3，
这个是钱包节点

3. 在中央节点上，用中央节点的地址，给WALLET_1，WALLET_2各发送10个币
这时候要带上 -mine 参数，意思是发送的时候立马挖一个新区块。
这是因为目前还没有 mine 节点。
然后启动这个中央节点 - `startnode`
这个服务会一直开着。

4. 再回到wallet node上，`startnode`，然后它会从中央节点那里下载所有区块
这时候可以通过查看余额的方式验证确实把数据都拿下来了。
* 在`startnode`之前，先create一个blockchain，但是不是用`createblockchain`的命令，而是复制 `cp blockchain_genesis.db blockchain_3001.db`
* 这时候可以检查 WALLET_1，WALLET_2的余额，应该都是 10。
同时，可以检查中央节点的钱包的余额，因为已经从中央节点把所有数据都下载下来了，应该也是10。

5. 再开一个终端，set NODE_ID = 3002，这是一个矿机节点
先生成一个钱包地址，复制blockchain `cp blockchain_genesis.db blockchain_3002.db`
然后 `startnode -miner MINER_WALLET`

6. 回到3001上，发送一些币给某些人
```
go run *.go send -from WALLET_1 -to WALLET_3 -amount 1
go run *.go send -from WALLET_2 -to MINER_WALLET -amount 1
```

7. 回到3002上，可以看到它挖了一个新区块

8. 回到3000上，可以看到，中央节点也收到了新区块

9. 回到wallet node，3001，`startnode`
它会下载所有新区块，停止服务，可以查看各个的余额。
```
go run *.go send -getbalance WALLET_1  // 9
go run *.go send -getbalance WALLET_2  // 9
go run *.go send -getbalance WALLET_3  // 1
go run *.go send -getbalance MINER_WALLET  // 11
```

OVER！

# 总结
作者没时间去实现一个完整的 P2P 网络了，但是期望能够通过前面的讲解回答读者的一些问题，并引起读者的思考。
比特币技术的背后还是有很多有意思的设计的。

此外，读者可以从实现一个 `addr` 消息着手进一步改进这个网络。
这个消息很重要，它能够让节点之间相互发现。

Links:
1. [Bitcoin protocol documentation](https://en.bitcoin.it/wiki/Protocol_documentation)
2. [Bitcoin network](https://en.bitcoin.it/wiki/Network)


