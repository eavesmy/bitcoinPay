package bcp

import (
	"bytes"
	"encoding/hex"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zdy23216340/bitcoinPay/lib"
	"image/png"
	"io/ioutil"
	"math/big"
)

type EthWallet struct {
	privateKey string
	address    string
}

func (w *EthWallet) New() WalletBase {

	key, _ := crypto.GenerateKey()
	w.privateKey = hex.EncodeToString(key.D.Bytes())
	w.address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	return w
}

func (w *EthWallet) Address() string {
	return w.address
}

func (w *EthWallet) QRCode() []byte {
	qrCode, _ := qr.Encode(w.address, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	fi := bytes.NewBuffer([]byte{})

	png.Encode(fi, qrCode)

	buf, _ := ioutil.ReadAll(fi)

	return buf
}

func (w *EthWallet) PrivateKey() string {
	return w.privateKey
}

func (w *EthWallet) Balance() *big.Int {
	balance := lib.GetEthBalance(w.address)
	return big.NewInt(balance)
}

func (w *EthWallet) BalanceOf(addr string) {

}

func (w *EthWallet) History() {

}

func (w *EthWallet) LastTransferIn()         {}
func (w *EthWallet) LastTransferOut()        {}
func (w *EthWallet) Transfer(addr string)    {}
func (w *EthWallet) QueryByTxid(txid string) {}
