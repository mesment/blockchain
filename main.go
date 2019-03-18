package main

import (
	"github.com/mesment/blockchain/block"
)

func main() {
	//blockchain := block.GetBlockchain()
	cli := block.CLI{}
	cli.Run()

	/*
	blockchain.AddBlockToBlockchain("new block")
	block := block.NewBlock("another block",blockchain)
	blockchain.AddBlockToBlockchain("another block")
	blockchain.PrintBlockchain()
	*/

}