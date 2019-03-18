package block

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	BC *Blockchain
}

func (cli *CLI) createBlockchainWithGenesisBlock(data string) {
	//CreateBlockchainWithGenesisBlock(data)
}

func (cli *CLI) Run(){
	cli.checkArgs();

	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	printBlockchainCmd := flag.NewFlagSet("printblockchain", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)

	addressData := createBlockchainCmd.String("address","","创建创世区块地址")
	sendBlockFrom := sendBlockCmd.String("from","","比特币发送者")
	sendBlockTo := sendBlockCmd.String("to","","比特币接收者")
	//sendBlockAmount := sendBlockCmd.String("amt","","发送金额")
	sendBlockAmount := sendBlockCmd.Int("amt",0,"发送金额")
	getBalanceAddressData := getBalanceCmd.String("address","","地址")

	switch os.Args[1] {

	case "printblockchain":
		err := printBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *addressData == "" {
			cli.printUsage()
			os.Exit(1)
		}
		cli.createBlockchain(*addressData)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddressData == ""  {
			cli.printUsage()
			os.Exit(1)
		}
		fmt.Printf("address:%s\n",*getBalanceAddressData)
		cli.getBalance(*getBalanceAddressData)
	}
	if sendBlockCmd.Parsed() {
		if *sendBlockFrom== "" || *sendBlockTo == "" || *sendBlockAmount == 0 {
			cli.printUsage()
			os.Exit(1)
		}
		fmt.Printf("from:%s\n",*sendBlockFrom)
		fmt.Printf("to:%s\n",*sendBlockTo)
		fmt.Printf("amt:%d\n",*sendBlockAmount)

		cli.send(*sendBlockFrom, *sendBlockTo,*sendBlockAmount)
	}
	if printBlockchainCmd.Parsed() {
		cli.printBlockchain()
	}

}

// 判断数据库是否存在
func DBExist()bool {
	if _,err :=os.Stat(DbName);os.IsNotExist(err) {
		return false
	}
	return true
}


func (cli *CLI)send(from string, to string, amt int){
	if !DBExist() {
		fmt.Printf("数据库不存在")
		os.Exit(1)
	}
	blockchain := GetBlockchain()
	defer blockchain.DB.Close()
	blockchain.MineNewBlock(from, to, amt)

}

func (cli *CLI)getBalance(address string) float64 {
	var balance  float64

	bc := GetBlockchain()
	defer bc.DB.Close()

	balance = bc.getBalance(address)
	return balance
}

func (cli *CLI)printBlockchain(){
	if !DBExist() {
		fmt.Printf("数据库不存在")
		os.Exit(1)
	}
	blockchain := GetBlockchain()
	defer blockchain.DB.Close()
	blockchain.PrintBlockchain()
/*
	for {
		block := iterator.Next()
		fmt.Printf("-------------------------------------Block Info Begin----------------------------------------\n")
		fmt.Printf("block height: %d\n", block.Height)
		fmt.Printf("block Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("block PreHash: %x\n", block.PreHash)
		fmt.Printf("block Data: %v\n", block.Transactions)
		fmt.Printf("block Hash: %x\n", block.Hash)
		fmt.Printf("block Nonce: %d\n", block.Nonce)
		fmt.Printf("-------------------------------------Block Info End------------------------------------------\n")
		//如果当前的区块是genesis则退出循环
		if block.PreHash == nil {
			break
		}
	}
*/
}

func (cli *CLI)printUsage(){
	fmt.Println("Usage:")
	fmt.Printf(" createblockchain -adderss \"交易数据\" \n")
	fmt.Printf(" printblockchain -打印区块链的所有区块信息\n")
	fmt.Printf(" send -from -to -amt  -从发送者给接受者转amt个比特币\n")
	fmt.Printf(" getbalance -address   -获取指定地址的比特币余额\n")
}

func (cli * CLI)createBlockchain(data string) {
	cli.BC = CreateBlockchainWithGenesisBlock(data)
	fmt.Println("blcockinfo :%v",cli.BC)
}

func(cli *CLI)checkArgs(){
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)

	}
}
