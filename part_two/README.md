# Part Two

[building-blockchain-in-go Part 2](https://jeiwan.cc/posts/building-blockchain-in-go-part-2/)

第一部分构造了一个简单的区块链模型，包括初始化一个区块链以及增加区块。
但是存在着很多问题，本章节将解决其中的一个问题：新增区块太容易的问题（Proof-of-Work）。

>In Bitcoin, the goal of such work is to find a hash for a block, that meets some requirements.
In the original Hashcash implementation, the requirement sounds like “first 20 bits of a hash must be zeros”.
In Bitcoin, the requirement is adjusted from time to time, because, by design, a block must be generated every 10 minutes,
despite computation power increasing with time and more and more miners joining the network.

比特币的要求就是会不断变化的，因为越来越多的人加入挖矿以及计算能力不断提升。
这里，我们的要求是，data+counter的哈希值的前24位是0。

注解：SHA256哈希后的值是一个256位的二进制数，转成16进制就是64位。
前24位是0的意思就是，哈希后的值**小于**
```
0x10000000000000000000000000000000000000000000000000000000000
即
0000010000000000000000000000000000000000000000000000000000000000
转成16进制，就是前6位是0
```

# Next
下一步将解决"钱包、地址、交易、Consensus"等问题。