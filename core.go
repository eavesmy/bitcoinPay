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
	BalanceOf(string, string) *big.Int
	History(...*Option) []*lib.Transaction
	LastTransferIn()
	Fee() string
	LastTransferOut()
	Transfer(string, string, ...map[string]string) error
	TokenTransfer(string, string, string, ...*Option) error
	QueryByTxid(string) *lib.Transaction
	ValidAddress(string) bool
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
	case "btc":
		wallet = &BtcWallet{privateKey: privateKey, id: 0}
	case "usdt":
		wallet = &BtcWallet{privateKey: privateKey, id: 31}
	}

	return wallet
}

type Option struct {
	From     string
	To       string
	Amount   string
	Data     string
	Gas      string
	Limit    int
	Nonce    string
	ChanId   int
	Contract string
    
    Page int
	Start   int
	End    int 
	Sort    string
	Address string
}

func (o *Option) Default() {
    if o.Sort == "" {
        o.Sort = "desc"
    }
}
