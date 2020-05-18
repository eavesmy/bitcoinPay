package bcp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/eavesmy/bitcoinPay/lib"
	"github.com/eavesmy/bitcoinPay/lib/eth"
	gtype "github.com/eavesmy/golang-lib/type"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"image/png"
	"io/ioutil"

	// "math"
	"math/big"
	"regexp"
	"strconv"
	// "sync"
)

const (
	GAS_LIMIT = 21000
	GAS_PRICE = 500000000000
)

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
func (w *EthWallet) History(options ...*Option) []*lib.Transaction {

	option := &Option{}

	if len(options) > 0 {
		option = options[0]
	}

	option.Default()

	if w.address == "" {
		fmt.Println("Error: param address missed")
		return []*lib.Transaction{}
	}

	return eth.GetTransactions(map[string]string{
		"start":   gtype.Int2String(option.Start),
		"end":     gtype.Int2String(option.End),
		"sort":    option.Sort,
		"address": w.address,
	})
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

// params: to,amount,contract
func (w *EthWallet) TokenTransfer(addr string, amount *big.Int, contract string, options ...*Option) (string, error) {

	option := &Option{}
	if len(options) > 0 {
		option = options[0]
	}

	option.From = w.address
	option.To = addr
	option.Amount = amount
	option.Data = "0xa9059cbb" + eth.StringPadding64(addr) + eth.StringPadding64(amount.String())
	option.Nonce = w.Nonce(w.address)

	hash := w.sign(option)

	return eth.SendRawTransaction(hash), nil
}

// Query transaction by txid
func (w *EthWallet) QueryByTxid(txid string) *lib.Transaction {
	return eth.GetTransactionByHash(txid)
}

// params: to:string privatekey:string amount:int data:[]byte gasLimit:int64 gasPrice:int64 chainid:int
func (w *EthWallet) sign(option *Option) string {

	fmt.Println("签名前参数检查: ", option)

	tx := types.NewTransaction(option.Nonce, common.HexToAddress(option.To), option.Amount, GAS_LIMIT, big.NewInt(GAS_PRICE), []byte(option.Data))

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

func (w *EthWallet) Nonce(addrs ...string) uint64 {
	addr := w.address
	if len(addrs) > 0 {
		addr = addrs[0]
	}
	nonce := eth.GetTransactionCount(addr)

	u, _ := strconv.ParseUint(nonce, 0, 64)

	return u
}
