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
	"regexp"
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
func (w *EthWallet) BalanceOf(addr string, contract string) *big.Int {
	balance := eth.GetTokenBalance(addr, contract)
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
func (w *EthWallet) Transfer(addr string, amount string, options ...map[string]string) error {

	// hex := sign(map[string]string{ })
	/*
			var option map[string]string
			if len(options) > 0 {
				option = options[0]
			}

			// sign prams:
				hex := w.sign(map[string]string{
					"from":       w.address,
					"privatekey": w.privateKey,
					"amount":     amount,
					"to":         addr,
				})
		// 处理 data

		ret := eth.SendRawTransaction(hex)

	*/
	return nil
}

func (w *BtcWallet) TokenTransfer(addr, amount, contract string, options ...*Option) error {

	/*
		option := &Option{}
		if len(options) > 0 {
			option = options[0]
		}

		option.From = w.address
		option.To = addr
		option.Amount = amount
		amount = eth.Str2Hex(amount)
		option.Data = "0xa9059cbb" + eth.StringPadding64(addr) + eth.StringPadding64(amount)

		hash := w.sign(option)

		fmt.Println("签名", hash)
				  return sign({address: from,privateKey,data: code,to: contract})
			 86         .then(hash => Eth.proxy.eth_sendRawTransaction(hash));
	*/

	return nil
}

// Query transaction by txid
func (w *EthWallet) QueryByTxid(txid string) *lib.Transaction {
	return eth.GetTransactionByHash(txid)
}

// params: to:string privatekey:string amount:int data:[]byte gasLimit:int64 gasPrice:int64 chainid:int
func (w *EthWallet) sign(option *Option) string {

	wg := &sync.WaitGroup{}
	wc := 0

	if option.Nonce == "" {
		wc++
		go func() {
			option.Nonce = eth.GetTransactionCount(option.From)
			wg.Done()
		}()
	}
	if option.Gas == "" {
		wc++
		go func() {
			option.Gas = eth.GasPrice()
			wg.Done()
		}()
	}

	wg.Add(wc)
	wg.Wait()

	u_nonce, _ := (&big.Int{}).SetString(option.Nonce[2:], 16)
	a_to := common.HexToAddress(option.To)
	i_amount, _ := strconv.ParseFloat(option.Amount, 64)
	u_amount := (&big.Int{}).SetUint64(math.Float64bits(i_amount))
	u_gasPrice, _ := (&big.Int{}).SetString(option.Gas[2:], 16)

	if option.Data == "" {
		option.Data = "0x"
	}

	b_data := []byte(option.Data)

	fmt.Println(u_nonce.Uint64())
	fmt.Println(u_gasPrice)
	fmt.Println(string(b_data))

	fmt.Println("检查参数", u_nonce.Uint64(), a_to, u_amount, u_limit, u_gasPrice, b_data)

	tx := types.NewTransaction(u_nonce.Uint64(), a_to, u_amount, u_limit, u_gasPrice, b_data)

	b_privatekey, _ := hexutil.Decode(w.privateKey)
	privatekey, _ := crypto.ToECDSA(b_privatekey)

	signed, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), privatekey)
	if err != nil {
		fmt.Println(err)
	}
	b, _ := rlp.EncodeToBytes(signed)
	return hex.EncodeToString(b)
}

func (w *EthWallet) Fee() string {
	return eth.GasPrice()
}

func (w *EthWallet) ValidAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}
