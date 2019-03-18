package block

//交易输入
type TXInput struct {
	Txid  		[]byte		//引用的交易ID
	Vout   		int			//引用的交易输出的索引
	ScriptSig  	string 		//钱包地址(锁定脚本)
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool{
	return in.ScriptSig == unlockingData
}


