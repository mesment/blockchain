package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

const subsidy = 10 //挖矿奖励

//交易结构
type Transaction struct {
	ID   []byte    //交易ID
	Vin  []TXInput //交易输入
	Vout []TXOut   //交易输出

}

//创世区块创建的Transaction
func NewCoinbaseTX(address string) *Transaction {

	txin := TXInput{[]byte{}, -1, "Genesis Data"}
	txout := TXOut{subsidy, address}
	tx := Transaction{nil, []TXInput{txin}, []TXOut{txout}}
	tx.SetID()
	return &tx
}

//判断交易是否是创世区块的交易
//创世区块的交易输入时空
func (tx *Transaction) IsCoinbaseTx() bool {
	return tx.ID == nil && len(tx.Vin) == 1 && tx.Vin[0].Vout == -1
}

func NewUTXOTransaction(from string, to string, amt int, bc *Blockchain) *Transaction {

	//存储需要花费的交易输出
	var unspendTxoutput = make(map[string][]int)
	//查找from的所有未输出交易
	utxos := bc.FindUnspendTransactions(from)

	var sumamt int

	//遍历所有未输出交易，找出金额大于等于amt的输出
Loop:
	for _, tx := range utxos {
		txID := hex.EncodeToString(tx.ID)
		for idx, out := range tx.Vout {
			if out.CanBeUnlockedWith(from) {
				if sumamt >= amt {
					break Loop
				} else {
					sumamt += out.Value
					unspendTxoutput[txID] = append(unspendTxoutput[txID], idx)
				}
			}
		}
	}
	if sumamt < amt {
		log.Panic("Avalible balance is not enough")
		return nil
	}
	//构造交易输入
	var vins []TXInput
	var vouts []TXOut

	for key, vouts := range unspendTxoutput {
		txid, _ := hex.DecodeString(key)

		for _,outidx := range vouts {
			vin := TXInput{txid, outidx, from}
			vins = append(vins, vin)
		}
	}

	//构造交易输出,判断是否需要找零
	vout := TXOut{amt, to}
	vouts = append(vouts, vout)
	if sumamt > amt {
		exchangeout := TXOut{sumamt - amt, from}
		vouts = append(vouts, exchangeout)
	}
	tx := Transaction{nil, vins, vouts}

	tx.SetID()

	return &tx
}

//设置交易ID
func (tx *Transaction) SetID() {
	var buf bytes.Buffer

	encode := gob.NewEncoder(&buf)
	err := encode.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(buf.Bytes())
	tx.ID = hash[:]
}
