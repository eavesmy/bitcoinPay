# 虚拟货币支付系统
开源的虚拟货币支付组件，基于 golang 开发

# 货币支持
- [ ] Eth
- [ ] Bitcoin
- [ ] Eos
- [ ] Xmr

# 示例
``` golang
package main

import(
    bcp "github.com/zdy23216340/bitcoinPay"
    "fmt"
)

const address = ""
const privateKey = ""

func main(){

	eth := bcp.Wallet("eth",privateKey)
	// or you can create a new wallet.
	// eth := bcp.Wallet("eth").New()

    b := eth.Balance() //BigNumber 

	params := map[string]string{}
	eth.Transfer("to",amount,params)

    transaction := eth.QueryByTxid(txid) // eth,btc,eos,xmr
    txs := eth.History()
	fmt.Println(transaction,txs)

	// check github.com/zdy23216340/bitcoinPay/core.go
}
```
