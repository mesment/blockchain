package main

import (
	"fmt"
	"github.com/mesment/blockchain/block"
)

func main() {
	//block := block.NewBlock("Genesis Blcok",1,nil)
	//genesisBlock := block.CreateGenesisBlock("Genesis Block")
	blockchain := block.CreateBlockchainWithGenesisBlock()
	lastblock := blockchain.Blocks[len(blockchain.Blocks) -1]
	newblock := block.NewBlock("new block",lastblock)
	blockchain.AddNewBlock(newblock)
	fmt.Println(blockchain)
	fmt.Println("block content:", newblock)
	blockbytes := newblock.SerializeBlock()
	fmt.Println("block bytes:", blockbytes)

	anotherBlock := block.DeSerializeBlock(blockbytes);
	fmt.Println("another block content:",anotherBlock)

}