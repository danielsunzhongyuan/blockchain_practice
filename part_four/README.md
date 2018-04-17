# Part Four
[building-blockchain-in-go Part 4 Transactions 1](https://jeiwan.cc/posts/building-blockchain-in-go-part-4/)

交易（transactions）是比特币的最核心的概念，而区块链的唯一目的就是为了用一种安全可靠的方式存储交易信息。
正因为交易非常重要和核心，所以这里将分两个章节介绍交易。
1. 本章节介绍交易的基本概念
2. 下一章节介绍更详细的内容

但是需要注意的是：
区块链里是一个公共开放的数据库，所以我们不存储关于钱包主人的敏感信息。
所以，没有账户、余额、地址、币、发送方和接收方的概念。

# 比特币的交易
交易包括3个部分：
- ID []byte
- Vin []TXInput
- Vout []TXOutput

一条交易的输入（Vin）来源与前一条交易的输出（这里涉及到第一条交易的问题，下面会详细讲解）。
而输出则真正存储了"币"的概念。注意：
1. 不是所有输出都会成为下一条交易的输入
2. 在一条交易里，输入可以来源于多条交易的输出
3. 一个输入必须来源与某个输出

*交易事实上只是通过脚本加上了锁的数据，只有加锁的人才能解锁*

# 交易的输出
```
type TXOutput struct {
    Value           int
    ScriptPubKey    string
}
```

上面即是输出的结构，其中Value就是用来存储"币"这一概念的，而ScriptPubKey则是用来加锁的。
>In Bitcoin, the value field stores the number of satoshis, not the number of BTC.
A satoshi is a hundred millionth of a bitcoin (0.00000001 BTC),
thus this is the smallest unit of currency in Bitcoin (like a cent).

目前还没有地址（addresses）的概念，所以先忽略script相关的逻辑。

ScriptPubKey会存储一个随机的字符串（即用户自定义的钱包地址）。
>By the way, having such scripting language means that Bitcoin can be used as a smart-contract platform as well.

需要注意的是：输出是不可见的，只会当交易成功后进行变化。


# 交易的输入
```
type TXInput struct {
    Txid        []byte
    Vout        int
    ScriptSig   string
}
```
- Txid 存储的是交易的ID
- Vout 存储的是交易里输出的索引
- ScriptSig 存储的是将被用在输出 ScriptPubKey 的数据，这个数据正确，输出才可以被解锁

每个交易，至少包括1个输入和1个输出

# 鸡生蛋 蛋生鸡问题
在比特币的概念里，是现有的Outputs，再有的Inputs
当一个矿工开始挖币的时候，它会增加一条 coinbase的交易信息。
Coinbase transaction是一种特殊类型的交易，它不要求有前一条的输出。
它直接生成Output。（所以说，Output先于Input出现）
这就是矿工挖币所获得的奖励。

在区块链概念里，最初始的区块称之为"Genesis block"，也就是这一个block，产生了区块链里的第一个Output。

# 在区块链里存储交易信息
每个区块都至少包含一个交易，否则无法继续挖矿。这样以来，Block里的data就将替换为transactions
```
type Block struct {
    Timestamp       int64
    Transactions    []*Transaction
    PrevBlockHash   []byte
    Hash            []byte
    Nonce           int
}
```


# 操作命令
```
➜  part_four git:(master) ✗ go run *.go -h
Usage:
  getbalance -address ADDRESS - Get balance of ADDRESS
  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS
  printchain - Print all the blocks of the blockchain
  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO
exit status 1


go run *.go createblockchain -address Ivan
go run *.go getbalance -address Ivan // Balance of 'Ivan': 10
go run *.go send -from Ivan -to John -amount 6
go run *.go getbalance -address Ivan // Balance of 'Ivan': 4
go run *.go getbalance -address John // Balance of 'John': 6

```