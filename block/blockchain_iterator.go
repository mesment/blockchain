package block

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	blockHash []byte	//当前节点的hash
	DB  *bolt.DB		//数据库
}

func (bci *BlockchainIterator)Next() *Block {
	var block *Block
	err := bci.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket != nil {
			blockBytes := bucket.Get(bci.blockHash)
			//获取当前迭代器中blockHash 代表的区块
			block = DeSerializeBlock(blockBytes)
			//更新迭代器中的区块hash
			bci.blockHash = block.PreHash
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return block
}