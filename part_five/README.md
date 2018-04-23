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
首先介绍钱包的概念：
```
type Wallet struct {
    PrivateKey ecdsa.PrivateKey
    PublicKey []byte
}

type Wallets struct {
    Wallets map[string]*Wallet
}
func NewWallet() *Wallet {
    private, public := newKeyPair()
    wallet := Wallet{private, public}
    return &wallet
}
func newKeyPair() (esdsa.PrivateKey, []byte) {
    curve := elliptic.P256()
    private, err := ecdsa.GenerateKey(curve, rand.Reader)
    pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
    return *private, pubKey
}
```

正如上面描述的那样，钱包其实就是一对秘钥。这对秘钥由ECDSA算法产生，其中私钥由elliptic curve生成，公钥由私钥生成。
公钥代表的是elliptic curve上的一个点，所以是包含的横坐标和纵坐标。
在比特币里，用横坐标+纵坐标表征公钥。

现在来生成地址：
```
func (w Wallet) GetAddress() []byte {
    pubKeyHash := HashPubKey(w.PublicKey)

    versionedPayload := append([]byte{version}, pubKeyHash...)
    checksum := checksum(versionedPayload)

    fullPayload := append(versionedPayload, checksum...)
    address := Base58Encode(fullPayload)

    return address
}

func HashPubKey(pubKey []byte) []byte {
    publicSHA256 := sha256.Sum256(pubKey)

    RIPEMD160Hasher := ripemd160.New()
    _, err := RIPEMD160Hasher.Write(publicSHA256[:])
    publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

    return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}
```

实现了地址之后，再来修正交易的输入和输出。
```
type TXInput struct {
    Txid        []byte
    Vout        int
    Signature   []byte
    PubKey      []byte
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
    lockingHash := HashPubKey(in.PubKey)
    return bytes.Compare(lockingHash, pubKeyHash) == 0
}

type TXOutput struct {
    Value       int
    PubKeyHash  []byte
}

func (out *TXOutput) Lock(address []byte) {
    pubKeyHash := Base58Decode(address)
    pubKeyHash := pubKeyHash[1 : len(pubKeyHash) - 4]
    out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
    return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
```

这里可以看到，我们不再使用 ScriptPubKey和ScriptSig，因为我们并不打算实现一个脚本语言。
取而代之的，是ScriptSig被分为两个部分：签名（Signature）和公钥（PubKey）。而且ScriptPubHash被重命名为PubKeyHash。

# 签名的实现
交易必须被签名，因为只有签名才能保证一个人不会去使用属于其他人的比特币。
如果签名无效，那么交易也就无效，也就无法被添加到区块链里去。

现在要考虑的事：要使用交易中的哪些数据来作为被签名的对象？是不是要把整个交易作为签名的对象呢？

需要签名的内容有：
1. Public key hashes stored in unlocked outputs. This identifies “sender” of a transaction.
2. Public key hashes stored in new, locked, outputs. This identifies “recipient” of a transaction.
3. Values of new outputs.

>比特币支持多种locking/unlocking逻辑，这些逻辑存储在ScriptSig和ScriptPubKey里。
所以比特币的签名对象里包含里ScriptPubKey的全部内容。

具体实现如下：
```
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
    if tx.IsCoinbase() {
        return
    }
    txCopy := tx.TrimmedCopy()

    for inID, vin := range txCopy.Vin {
        prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
        txCopy.Vin[inID].Signature = nil
        txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
        txCopy.ID = txCopy.Hash()
        txCopy.Vin[inID].PubKey = nil

        r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
        signature := append(r.Bytes(), s.Bytes()...)
        tx.Vin[inID].Signature = signature
    }
}

func (tx *Transaction) TrimmedCopy() Transaction {
    var inputs []TXInput
    var outputs []TXOutput

    for _, vin := range tx.Vin {
        intputs = append(intputs, TXInput{vin.Txid, vin.Vout, nil, nil})
    }

    for _, vout := range tx.Vout {
        outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
    }

    txCopy := Transaction{tx.ID, inputs, outputs)

    return txCopy
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}
```

上面的代码包括了签名和校验。

现在需要重新定义 FindTransaction/SignTransaction/VerifyTransaction 等。
```
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}
```

对交易的签名发生在 `NewUTXOTransaction` 里。
```
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	...

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}
```

而对交易的校验则发生在把交易放进区块的时候。
```
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}
	...
}
```


# 总结
目前为止，我们已经实现了比特币相关的，除了网络之外的大部分特性，下一章节，我们会结束交易的讲解。
