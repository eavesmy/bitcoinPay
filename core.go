package bcp

import (
	"math/big"
	"strings"
)

type WalletBase interface {
	New() WalletBase
	QRCode() []byte
	Address() string
	PrivateKey() string
	Balance() *big.Int
	BalanceOf(string)
	History()
	LastTransferIn()
	LastTransferOut()
	Transfer(string)
	// QueryByTxid(string, string)
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
