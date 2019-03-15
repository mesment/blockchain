package main

import (
	"github.com/mesment/blockchain/block"
)

func main() {
	blockchain := block.GetBlockchain()
	defer blockchain.DB.Close()
	cli := block.CLI{blockchain}
	cli.Run()

	/*
	blockchain.AddBlockToBlockchain("new block")
	block := block.NewBlock("another block",blockchain)
	blockchain.AddBlockToBlockchain("another block")
	blockchain.PrintBlockchain()
	*/

}