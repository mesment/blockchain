package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Height 		int64		//区块高度
	Timestamp 	int64		//时间戳
	PreHash		[]byte		//前一个区块的hash
	Data 		[]byte		//交易数据
	Hash 		[]byte		//区块的hash值
	Nonce		int64 		//随机数
}


func NewBlock(data string,prevBlock *Block) *Block {
	var block = Block{}
	block.Height = prevBlock.Height + 1
	block.PreHash = prevBlock.Hash
	block.Data = []byte(data)
	block.Nonce = 0

	//挖矿，通过工作量证明计算Nonce和hash值
	//使得当前区块的hash值开头0的个数与难度系数相同
	for {
		//更新区块的时间戳
		block.Timestamp= time.Now().Unix()
		blockHashStr := hex.EncodeToString(block.CalculateHash(block.Nonce))
		fmt.Println("挖矿中 ",blockHashStr)

		//检查hash是否满足难度
		if IsHashValid(blockHashStr) {
			if VerifyBlock(&block,prevBlock) {
				fmt.Println("挖矿成功！")
				block.Hash = []byte(blockHashStr)
				break
			}
		}

		block.Nonce++
	}
	return &block
}

//计算区块的hash值
func (b *Block)CalculateHash(nonce int64) []byte {
	blockInfoStr := strconv.FormatInt(b.Height,10)+ strconv.FormatInt(b.Timestamp,10) +
				hex.EncodeToString(b.PreHash) + hex.EncodeToString(b.Data)+ strconv.FormatInt(nonce,10)
	h := sha256.New()
	h.Write([]byte(blockInfoStr))
	hashed := h.Sum(nil)
	b.Hash = hashed
	return hashed
}

//创建创世区块
func CreateGenesisBlock(data string)*Block {
	genesis :=Block{1, time.Now().Unix(), nil, []byte(data),nil,0}
	genesis.CalculateHash(0)
	return &genesis
}

//判断hash是否有效，hash值是否满足难度
func IsHashValid(hashStr string) bool {
	difficutStr := strings.Repeat("0",4)
	fmt.Printf("Difficult string is %s:",difficutStr)

	if strings.HasPrefix(hashStr, difficutStr) {
		return true
	}
	return false
}

//校验区块是否合法
func VerifyBlock(newBlock, prevBlock *Block) bool {
	if newBlock.Height != prevBlock.Height + 1 {
		return false
	}
	if string(newBlock.PreHash) != string(prevBlock.Hash) {
		return false
	}
	return true
}

//将区块block序列化成字节数组
func (block *Block)SerializeBlock() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//将字节数组反序列化成Block
func DeSerializeBlock(blockBytes []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))

	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return &block
}


