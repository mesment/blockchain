package block

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type CLI struct {
	BC *Blockchain
}

func (cli *CLI) createBlockchainWithGenesisBlock(data string) {
	//CreateBlockchainWithGenesisBlock(data)
}

func (cli *CLI) Run(){
	cli.checkArgs();

	addBlockCmd := flag.NewFlagSet("addblock",flag.ExitOnError)
	printBlockchainCmd := flag.NewFlagSet("printblockchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data","","交易数据")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printblockchain":
		err := printBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}
	if printBlockchainCmd.Parsed() {
		cli.printBlockchain()
	}

}

func (cli *CLI)printBlockchain(){
	iterator := cli.BC.Iterator()

	for {
		block := iterator.Next()
		fmt.Printf("-------------------------------------Block Info Begin----------------------------------------\n")
		fmt.Printf("block height: %d\n", block.Height)
		fmt.Printf("block Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("block PreHash: %x\n", block.PreHash)
		fmt.Printf("block Data: %s\n", block.Data)
		fmt.Printf("block Hash: %x\n", block.Hash)
		fmt.Printf("block Nonce: %d\n", block.Nonce)
		fmt.Printf("-------------------------------------Block Info End------------------------------------------\n")
		//如果当前的区块是genesis则退出循环
		if block.PreHash == nil {
			break
		}
	}
}

func (cli *CLI)printUsage(){
	fmt.Println("Usage:")
	fmt.Printf("./appname addblock -data \"交易数据\" -添加一个新的区块到区块链\n")
	fmt.Printf("./appname printblockchain -打印区块链的所有区块信息\n")
}

func (cli * CLI)addBlock(data string) {
	cli.BC.AddBlockToBlockchain(data)
}

func(cli *CLI)checkArgs(){
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)

	}
}
