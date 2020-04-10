package bcp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/eavesmy/bitcoinPay/lib"
	"github.com/eavesmy/bitcoinPay/lib/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"image/png"
	"io/ioutil"
	"math"
	"math/big"
	"strconv"
	"sync"
)

const u_limit = 21000

// Based on WalletBase
type EthWallet struct {
	privateKey string
	address    string
}

// Create New Wallet.
func (w *EthWallet) New() WalletBase {

	key, _ := crypto.GenerateKey()
	w.privateKey = hex.EncodeToString(key.D.Bytes())
	w.address = crypto.PubkeyToAddress(key.PublicKey).Hex()

	return w
}

// TODO: GetAddress method
// Through private key generate address.

// Get wallet address or set address to wallet.
func (w *EthWallet) Address(addrs ...string) string {
	if len(addrs) > 0 {
		w.address = addrs[0]
	}
	return w.address
}

// Get QR of address.
func (w *EthWallet) QRCode() []byte { // for receive coin
	qrCode, _ := qr.Encode(w.address, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	fi := bytes.NewBuffer([]byte{})

	png.Encode(fi, qrCode)

	buf, _ := ioutil.ReadAll(fi)

	return buf
}

// Get address private key.
func (w *EthWallet) PrivateKey() string {
	return w.privateKey
}

// Get balance of address. Unit: wei
func (w *EthWallet) Balance() *big.Int {
	balance := eth.GetEthBalance(w.address)
	i_balance, _ := strconv.Atoi(balance)
	return big.NewInt(int64(i_balance))
}

// Like Balance().
func (w *EthWallet) BalanceOf(addr string) *big.Int {
	balance := eth.GetEthBalance(addr)
	i_balance, _ := strconv.Atoi(balance)
	return big.NewInt(int64(i_balance))
}

// Get history transactions
func (w *EthWallet) History(params ...map[string]string) []*lib.Transaction {
	param := map[string]string{}

	if len(params) > 0 {
		param = params[0]
	}

	_, exists := param["start"]
	if !exists {
		param["start"] = "0"
	}

	_, exists = param["end"]
	if !exists {
		param["end"] = "latest"
	}

	_, exists = param["sort"]
	if !exists {
		param["sort"] = "asc"
	}

	_, exists = param["address"]
	if !exists {
		fmt.Println("Error: param address missed")
		return []*lib.Transaction{}
	}

	return eth.GetTransactions(param)

}

// Get latest transfaction with coin in
func (w *EthWallet) LastTransferIn() {}

// Get lastest transfaction with coin out
func (w *EthWallet) LastTransferOut() {}

// Transfer
// Local sign and use etherscan api to transfer.
func (w *EthWallet) Transfer(addr string, amount float64, options map[string]string) {

	// sign prams:
	// hex := sign(map[string]string{ })

	hex := w.sign(map[string]string{
		"from":       w.address,
		"privatekey": w.privateKey,
		"amount":     strconv.FormatFloat(amount, 'f', 6, 64),
		"to":         addr,
	})

	ret := eth.SendRawTransaction(hex)
	fmt.Println(ret)
    
    // TODO:
	// test not finish.
}

// Query transaction by txid
func (w *EthWallet) QueryByTxid(txid string) *lib.Transaction {
	return eth.GetTransactionByHash(txid)
}

// params: to:string privatekey:string amount:int data:[]byte gasLimit:int64 gasPrice:int64 chainid:int
func (w *EthWallet) sign(params map[string]string) string {

	wg := &sync.WaitGroup{}
	wc := 0

	nonce, exists := params["nonce"]
	if !exists {
		wc++
		go func() {
			nonce = eth.GetTransactionCount(params["from"])
			wg.Done()
		}()
	}
	gasPrice, exists := params["gasPrice"]
	if !exists {
		wc++
		go func() {
			gasPrice = eth.GasPrice()
			wg.Done()
		}()
	}

	wg.Add(wc)
	wg.Wait()

	u_nonce, _ := (&big.Int{}).SetString(nonce[2:], 16)
	a_to := common.HexToAddress(params["to"])
	i_amount, _ := strconv.ParseFloat(params["amount"], 64)
	u_amount := (&big.Int{}).SetUint64(math.Float64bits(i_amount))
	u_gasPrice, _ := (&big.Int{}).SetString(gasPrice[2:], 16)
	b_data := []byte{}

	data, exists := params["data"]
	if exists {
		b_data = []byte(data)
	} else {
		b_data = nil
	}

	tx := types.NewTransaction(u_nonce.Uint64(), a_to, u_amount, u_limit, u_gasPrice, b_data)

	b_privatekey, _ := hexutil.Decode(params["privatekey"])
	privatekey, _ := crypto.ToECDSA(b_privatekey)

	signed, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), privatekey)
	if err != nil {
		fmt.Println(err)
	}
	b, _ := rlp.EncodeToBytes(signed)
	return hex.EncodeToString(b)
}
