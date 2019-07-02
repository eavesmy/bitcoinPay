package bcp

func NewBtcWallet(pks ...string) {
	if len(pks) > 0 {

	}
}

type BtcWallet struct {
	PrivateKey string
}

func (w *BtcWallet) Balance() {

}

func (w *BtcWallet) BalanceOf(addr string) {

}

func (w *BtcWallet) History() {

}

func (w *BtcWallet) LastTransferIn()         {}
func (w *BtcWallet) LastTransferOut()        {}
func (w *BtcWallet) Transfer(addr string)    {}
func (w *BtcWallet) QueryByTxid(txid string) {}
