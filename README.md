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

function main(){
    eth := bcp.NewEthWallet(privateKey)
    btc := bcp.NewBitcoinWallet()
    eos := bcp.NewEosWallet()

    eth.Balance() //BigNumber 
    b := btc.Balance()
    b = eos.Balance()

    fmt.Println(b.ToString())

    b.History()
    b.LastTransferIn()
    b.LastTransferOut()
    
    txid := eth.Transfer("address",bcp.BigNumber(0.0001),"data")

    fmt.Println(txid)

    transactions := bcp.QueryByTxid(txid) // eth,btc,eos,xmr

    fmt.Println(transactions)
}
```
