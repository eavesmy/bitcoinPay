package bcp

import (
	"github.com/eavesmy/bitcoinPay/lib"
	"math/big"
	"strings"
)

type WalletBase interface {
	New() WalletBase
	QRCode() []byte
	Address(...string) string
	PrivateKey() string
	Balance() *big.Int
	BalanceOf(string) *big.Int
	History(...map[string]string) []*lib.Transaction
	LastTransferIn()
	LastTransferOut()
	Transfer(string, float64, map[string]string)
	QueryByTxid(string) *lib.Transaction
}

func Wallet(chain string, pks ...string) WalletBase {

	chain = strings.ToLower(chain)

	var wallet WalletBase

	privateKey := ""

	if len(pks) > 0 {
		privateKey = pks[0]
	}

	switch chain {
	case "eth":
		wallet = &EthWallet{privateKey: privateKey}
	}

	return wallet
}
