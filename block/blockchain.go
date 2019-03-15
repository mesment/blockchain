package block

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

type Blockchain struct {
	LastBlockHash []byte 	//最后一个区块的hash值
	DB  *bolt.DB			//区块数据库
}

const DbName = "blockchain.db"			//数据库文件名
const BucketName = "BlockBucket"		//boltdb 数据库存储bucket名
const LastBlockKey ="LastBlockHash"		//数据库存储的最后一个区块的hash键

//定义挖矿难度系数
const difficuty = 2


func GetBlockchain() *Blockchain {
	db,err := bolt.Open(DbName,0600,nil)
	var blockHash []byte
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket != nil {
			blockHash = bucket.Get([]byte(LastBlockKey))
		} else { // bucket is nil
			fmt.Println("Block chain is not exist,creating one")
			genesis :=CreateGenesisBlock("Genesis Block")
			blockBytes := genesis.SerializeBlock()
			blockHash = genesis.Hash

			bucket,err := tx.CreateBucket([]byte(BucketName))
			if err != nil {
				log.Panic(err)
			}

			err = bucket.Put(blockHash,blockBytes)
			if err != nil {
				log.Panic(err)
			}

			//存储最后一个区块的hash
			err = bucket.Put([]byte(LastBlockKey),blockHash)
			if err !=nil {
				log.Panic(err)
			}
		}
		return nil
	})

	return &Blockchain{blockHash,db}
}
func CreateBlockchainWithGenesisBlock() *Blockchain  {
	//创建或者打开数据库
	db,err := bolt.Open(DbName,0600,nil)
	if err != nil {
		fmt.Errorf("Open database %s failed,%v\n",DbName,err)
	}
	//defer db.Close()

	var blockHash []byte

	err = db.Update( func(tx *bolt.Tx) error {
		bucket,err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			log.Fatal("Create bucket %s failed,%v",BucketName,err)
		}

		genesis :=CreateGenesisBlock("Genesis Block")
		blockBytes := genesis.SerializeBlock()
		blockHash = genesis.CalculateHash(0)

		err = bucket.Put(blockHash,blockBytes)
		if err != nil {
			log.Panic(err)
		}
		//存储最后一个区块的hash
		err = bucket.Put([]byte(LastBlockKey),genesis.Hash)
		if err !=nil {
			log.Panic(err)
		}

		return nil
	})

	return &Blockchain{blockHash,db}
}


func (bc *Blockchain)AddBlockToBlockchain(data string) error {
	var newBlock *Block = NewBlock(data,bc)

	err := bc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket != nil {
			//将区块序列化并存储到数据库
			blockBytes := newBlock.SerializeBlock()
			err := bucket.Put(newBlock.Hash,blockBytes)
			if err != nil {
				return fmt.Errorf("Add new block to blockchain failed,%v",err)
			}

			//更新数据库中最后一个区块的索引
			err = bucket.Put([]byte(LastBlockKey),newBlock.Hash)
			if err != nil {
				return  fmt.Errorf("Update blockchain last block hash failed,%v",err)
			}
			//更新区块链的最后一个区块索引
			bc.LastBlockHash = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain)Iterator() *BlockchainIterator {

	return &BlockchainIterator{bc.LastBlockHash,bc.DB}
}

func (bc *Blockchain) PrintBlockchain() {

	var block *Block
	//通过迭代器遍历区块链
	iterator := bc.Iterator()
	for {
		block = iterator.Next()
		fmt.Printf("-------------------------------------Block Info Begin---------------------------------------------------\n")
		fmt.Printf("block height: %d\n", block.Height)
		fmt.Printf("block Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("block PreHash: %x\n", block.PreHash)
		fmt.Printf("block Data: %s\n", block.Data)
		fmt.Printf("block Hash: %x\n", block.Hash)
		fmt.Printf("block Nonce: %d\n", block.Nonce)
		fmt.Printf("-------------------------------------Block Info End---------------------------------------------------\n")
		//如果当前的区块是genesis则退出循环
		if block.PreHash == nil {
			break
		}
	}



}


