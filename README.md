# 虚拟货币支付系统
开源的虚拟货币支付组件，基于 golang 开发

# 货币支持
- [ ] Eth
- [ ] Bitcoin
- [ ] Eos
- [ ] Xmr

# Api
### 创建钱包
``` golang
package main

import(
    bcp "github.com/zdy23216340/bitcoinPay"
    "fmt"
)

func main(){

	eth := bcp.Wallet("eth",privateKey)
	// or you can create a new wallet.
	// eth := bcp.Wallet("eth").New()

    b := eth.Balance() //BigNumber 

    fmt.Println(b.ToString())

    eth.History()
    eth.LastTransferIn()
    eth.LastTransferOut()
    
    // txid := eth.Transfer(otherEthWalletAddress,bcp.BigNumber(0.0001),"data")

    //transactions := eth.QueryByTxid(txid) // eth,btc,eos,xmr

    // fmt.Println(transactions)
}
```
