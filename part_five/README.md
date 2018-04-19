# Part Five
[building-blockchain-in-go Part 5 Addresses](https://jeiwan.cc/posts/building-blockchain-in-go-part-5/)

上一章介绍的transaction的基本概念，并简单提到了address的概念，
本章将详细介绍address等概念。

# 比特币地址（Bitcoin address）
第一个比特币地址属于"Satoshi Nakamoto" - `1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa`。

比特币地址是公开的，所以你要知道对方的地址才能给他发送。

地址是唯一的，但还不能表征每个人的钱包（wallet）。比特币里，每个人的唯一标识是1个或多个公钥私钥对（a pair or pairs of private and public keys）。
比特币依赖于加密算法来生成这些key，并保证别人无法破解。

# Public-key Cryptography
一个公钥，一个私钥，私钥才是代表用户的唯一标识

# Digital Signatures 电子签名
电子签名的意思是，算法保证：
1. 数据从一方传到另一方时不会被改变
2. 数据是由具体某个发送方产生的
3. 发送方也无法否认数据已经被传输了

*用私钥进行签名，用公钥进行验证*

签名的时候需要包括：
1. 要被签名的数据
2. 私钥

验证的时候需要包括：
1. 被签名的数据
2. 签名
3. 公钥

再回到比特币里的交易，每个交易在生成的时候都要进行签名，然后才能放入区块中。
而验证是指：
1. 检查交易输入是否具有使用上一个交易输出的权限
2. 检查交易的签名是否正确


现在再重新回顾交易的整个生命周期，也就是比特币交易的流程：
1. 首先，生成一个"genesis block"，包含了一个 "coinbase transaction"
对于coinbase transaction而言，并没有交易的输入，所以签名是没必要的。
coinbase transaction的输出包含了一个哈希后的公钥
2. 然后发送比特币，交易产生。交易输入来自上游交易的输出。所有输入都包含一个未哈希的公钥，和整个交易的签名
3. 比特币网络的其他节点会验证它们接收到交易，也会检查交易输入里的公钥的哈希值（以此保证发送方发送的是它自己的比特币）；
签名是正确的（以此保证这条交易是发送方本人发送的）
4. 当一个挖矿节点准备好挖一个新的区块的时候，它会把交易放在区块里，然后再挖
5. 当新区块被挖到后，所有的其他矿机都会被通知到，并把这个区块放到区块链里
6. 新区块加入到区块链后，交易就完成了，这里的交易输出就可以作为下一个交易的基础了

# 椭圆曲线加密算法（elliptic curve cryptography）
如何保证公钥/私钥的唯一性呢？

比特币采用了椭圆曲线的方式来生成私钥。这种曲线将从 0～2\*\*（2\*\*56）中挑选一个数字。2\*\*（2\*\*56）大约是10**77。
同时，比特币使用ECDSA（Elliptic Curve Digital Signature Algorithm）算法来对交易进行签名。（我们也会用到这个算法）

# Base58
类似Base64，只是少了 0，O，I，l，+，/

还原后的公钥包括3个部分：Version，Public key hash，Checksum

下面就是一些基本代码实现了

# 地址的实现
