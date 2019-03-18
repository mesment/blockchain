package block

import (
	"encoding/hex"
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
			genesis :=CreateGenesisBlock([]*Transaction{&Transaction{}})
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
func CreateBlockchainWithGenesisBlock(address string) *Blockchain  {
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
		transaction := NewCoinbaseTX(address)
		genesis :=CreateGenesisBlock([]*Transaction{transaction})
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

func (bc *Blockchain)MineNewBlock(from string, to string, amt int) {

	tx := NewUTXOTransaction(from,to,amt,bc)
	txs :=[]*Transaction{tx}
	newblock := bc.AddBlockToBlockchain(txs)
	fmt.Printf("newblock:%v\n",newblock)


}

func (bc *Blockchain)getBalance(address string) float64 {
	var balance  float64
	utxo := bc.FindUTXO(address)
	if utxo != nil {
		for _, v := range utxo {
			balance += float64(v.Value)
		}
	}
	fmt.Printf("Balance of %s is:%f\n",address,balance)
	return balance
}

//遍历区块链bc，返回属于address的所有未花费交易
func (bc *Blockchain)FindUnspendTransactions(address string) []Transaction {
	//已花费交易输出map，键为已花费交易id，值为交易的输出索引数组
	var spendTxs map[string] []int = make(map[string][]int, 0)
	var unspendTxs []Transaction  //未花费交易

	iterator := bc.Iterator()
	for {
		block := iterator.Next()
		for _,tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
			Outputs:
				for idx, out := range tx.Vout {
					//判断交易输出out是否已花费
					if spendTxs[txID] != nil {
						for _,spendoutIdx := range spendTxs[txID] {
							if spendoutIdx == idx {
								continue Outputs
							}
						}

					}
					//out未花费，继续判断out是否属于address
					if out.CanBeUnlockedWith(address) == true {
						unspendTxs = append(unspendTxs,*tx)
					}
				}

				//判断是否是创世区块的交易
				if tx.IsCoinbaseTx() == false {
					for _, in := range tx.Vin {
						//判断交易输入in是否属于address
						if in.CanUnlockOutputWith(address) {
							inTxID := hex.EncodeToString(in.Txid)
							spendTxs[inTxID] = append(spendTxs[inTxID],in.Vout)
						}
					}
				}

		}
		//如果是创世区块则退出循环
		if block.PreHash == nil {
			break;
		}

	}
	return unspendTxs
}

func(bc *Blockchain)FindUTXO(address string) []TXOut{
	var UTXOs []TXOut
	unspendTXs := bc.FindUnspendTransactions(address)

	for _, tx := range unspendTXs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs,out)
			}
		}
	}

	return UTXOs
}

//将交易txs 打包并添加到区块链
func (bc *Blockchain)AddBlockToBlockchain(txs []*Transaction) error {
	var newBlock *Block = NewBlock(txs,bc)

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
		fmt.Printf("block Data: %v\n", block.Transactions)
		fmt.Printf("block Hash: %x\n", block.Hash)
		fmt.Printf("block Nonce: %d\n", block.Nonce)
		for _,tx := range block.Transactions {
			fmt.Printf("%x\n",tx.ID)
			fmt.Printf("vin:\n")
			for _,in := range tx.Vin {
				fmt.Printf("vinid:%x\n",in.Txid)
				fmt.Printf("vinout:%x\n",in.Vout)
				fmt.Printf("vinaddress:%x\n",in.ScriptSig)
			}
			fmt.Printf("vout:\n")
			for _,out := range tx.Vout {
				fmt.Printf("%v\n",out.Value)
				fmt.Printf("%v\n",out.ScriptPubKey)
			}
		}
		fmt.Printf("-------------------------------------Block Info End---------------------------------------------------\n")
		//如果当前的区块是genesis则退出循环
		if block.PreHash == nil {
			break
		}
	}



}


