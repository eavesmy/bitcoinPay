package bcp

func NewEosWallet(pks ...string) {
	if len(pks) > 0 {

	}
}

type EosWallet struct {
	PrivateKey string
}

func (w *EosWallet) Balance() {

}

func (w *EosWallet) BalanceOf(addr string) {

}

func (w *EosWallet) History() {

}

func (w *EosWallet) LastTransferIn()         {}
func (w *EosWallet) LastTransferOut()        {}
func (w *EosWallet) Transfer(addr string)    {}
func (w *EosWallet) QueryByTxid(txid string) {}
