# Part One

[building-blockchain-in-go Part 1](https://jeiwan.cc/posts/building-blockchain-in-go-part-1/)

这只是一个简单的区块链的模型：只是一些区块的集合，每个区块具有和上一个区块的连接。

实际的区块链要复杂的多，例如：新增一个区块需要通过完成非常高难度的计算来证明你有这样的能力（Proof-of-Work）；

同时，区块链是一个分布式的数据库，并没有中心决策节点（no single decision maker）。

因此，每一个新增的区块，都需要得到确认以及网络中其他参与者的同意（consensus）。

而且，我们的区块链里还没有交易。

# How to run the code
```
➜  part_one git:(master) ✗ go run *.go
Prev. hash:
Data: Genesis Block
Hash: 3b8fb06431fe1f18aa1ff4b7fb9104bc18d56cbeda7178926cb3115ea4d2a93e

Prev. hash: 3b8fb06431fe1f18aa1ff4b7fb9104bc18d56cbeda7178926cb3115ea4d2a93e
Data: Send 1 BTC to Ivan
Hash: 1d91dd3a96c4f0aeda38a7f8806b66c113c0b75ba6f584a28378e8b0959d3cac

Prev. hash: 1d91dd3a96c4f0aeda38a7f8806b66c113c0b75ba6f584a28378e8b0959d3cac
Data: Send 2 more BTC to Ivan
Hash: 6a9d98cc98f53f1bf02f3f26ac628e28f36d71d5e81eefb80409094dca8d387d

或者
➜  part_one git:(master) ✗ go test
```

# Next
下一章节将解决"添加区块太容易"的问题。