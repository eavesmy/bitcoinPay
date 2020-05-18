package bcp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	// "github.com/btcsuite/btcd/btcjson"
	// wbtc "github.com/dabankio/wallet-core/core/btc"

	// "github.com/dabankio/wallet-core/core/omni"
	"image/png"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/eavesmy/bitcoinPay/lib"
	"github.com/eavesmy/bitcoinPay/lib/btc"
	gtype "github.com/eavesmy/golang-lib/type"
)

type BtcWallet struct {
	privateKey string
	address    string
	id         int
}

func (w *BtcWallet) New() WalletBase {

	priKey, pubKey := btc.NewPair()

	w.address = btc.GetAddress(pubKey)
	w.privateKey = hex.EncodeToString(priKey)

	return w
}

// Get wallet address or set address to wallet.
func (w *BtcWallet) Address(addrs ...string) string {
	if len(addrs) > 0 {
		w.address = addrs[0]
	}
	return w.address
}

func (w *BtcWallet) QRCode() []byte { // for receive coin
	qrCode, _ := qr.Encode(w.address, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	fi := bytes.NewBuffer([]byte{})

	png.Encode(fi, qrCode)

	buf, _ := ioutil.ReadAll(fi)

	return buf
}

func (w *BtcWallet) PrivateKey() string {
	return w.privateKey
}

func (w *BtcWallet) Balance() *big.Int {
	amount := btc.GetBalance(w.address, w.id)
	return big.NewInt(amount)
}

func (w *BtcWallet) Balances(addrs []string, ids ...string) (ret map[string]string) {
	if len(addrs) == 0 {
		return
	}
	id := "0"
	if len(ids) > 0 {
		id = ids[0]
	}
	return btc.GetBalances(addrs, id)
}

func (w *BtcWallet) BalanceOf(addr string, id string) *big.Int {
	return nil
}

func (w *BtcWallet) History(options ...*Option) []*lib.Transaction {
	option := &Option{}

	if len(options) > 0 {
		option = options[0]
	}
	option.Default()

	return btc.GetHistory(w.address, option.Page)
}

// 获取手续费
func (w *BtcWallet) Fee() string {
	// 获取最后一笔账单
	return btc.GetFee()
}

func (w *BtcWallet) LastTransferIn()  {}
func (w *BtcWallet) LastTransferOut() {}
func (w *BtcWallet) Transfer(addr string, amount string, options ...map[string]string) error {

	var option map[string]string
	if len(options) > 0 {
		option = options[0]
	}

	fmt.Println(option)

	return nil
}

type ListUnspentResult struct {
	TxID string `json:"txid"`
}

func (w *BtcWallet) TokenTransfer(addr string, amount *big.Int, contract string, options ...*Option) (string, error) {

	// unspent := new(wbtc.BTCUnspent)
	// unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, utxo.RedeemScript)

	return "", nil
}

func (w *BtcWallet) QueryByTxid(txid string) *lib.Transaction { return nil }

func (w *BtcWallet) ValidAddress(address string) bool {

	len := len(address)
	if len < 25 {
		return false
	}

	if strings.HasPrefix(address, "1") {
		if len >= 26 && len <= 34 {
			return true
		}
	}

	if strings.HasPrefix(address, "3") && len == 34 {
		return true
	}

	if strings.HasPrefix(address, "bc1") && len > 34 {
		return true
	}

	return false
}

func (w *BtcWallet) Nonce(addrs ...string) uint64 {
	addr := w.address
	if len(addrs) > 0 {
		addr = addrs[0]
	}
	nonce := btc.GetTransactionCount(addr)
	i_nonce := gtype.String2Int(nonce)
	return uint64(i_nonce)
}

func (w *BtcWallet) sign(option *Option) string {
	return ""
}
