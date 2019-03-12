package block

type Blockchain struct {
	Blocks []*Block
}

//定义挖矿难度系数
const difficuty = 3

func CreateBlockchainWithGenesisBlock() *Blockchain {
	genesis :=CreateGenesisBlock("Genesis Block")
	return &Blockchain{[]*Block{genesis}}
}

func (bc *Blockchain)AddNewBlock(newBlock *Block){
	bc.Blocks = append(bc.Blocks,newBlock)
	return
}