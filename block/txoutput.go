package block


//交易输出
type TXOut struct{
	Value 			int     //金额 以分为单位
	ScriptPubKey	string	//钱包地址(解锁脚本)
}

func (out *TXOut)CanBeUnlockedWith(unlockingData string) bool{
	return out.ScriptPubKey == unlockingData
}