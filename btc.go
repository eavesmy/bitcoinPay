package bcp

import (
	"bytes"
    "strings"
	"encoding/hex"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/eavesmy/bitcoinPay/lib"
	"github.com/eavesmy/bitcoinPay/lib/btc"
	"image/png"
	"io/ioutil"
	"math/big"
)

type BtcWallet struct {
	privateKey string
	address    string
    id int
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
    amount := btc.GetBalance(w.address,w.id)
    return big.NewInt(amount)
}

func (w *BtcWallet) Balances(addrs []string,ids ...string) (ret map[string]string){
    if len(addrs) == 0 {
        return
    }
    id := "0"
    if len(ids) > 0 {
        id = ids[0]
    }
    return btc.GetBalances(addrs,id)
}

func (w *BtcWallet) BalanceOf(addr string,id string) *big.Int {
	return nil
}

func (w *BtcWallet) History(params ...map[string]string) []*lib.Transaction {
	return nil
}

// 获取手续费
func (w *BtcWallet) Fee() string{
    // 获取最后一笔账单
    return btc.GetFee()
}

func (w *BtcWallet) LastTransferIn()                                                 {}
func (w *BtcWallet) LastTransferOut()                                                {}
func (w *BtcWallet) Transfer(addr string, amount float64, options map[string]string) {}
func (w *BtcWallet) QueryByTxid(txid string) *lib.Transaction                        { return nil }

func (w *BtcWallet) sign(param map[string]string) string {
	return ""
}

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
