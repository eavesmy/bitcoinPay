package btc

/*
{
      "divisible": true,
      "error": true,
      "id": "0",
      "pendingneg": "0",
      "pendingpos": "0",
      "propertyinfo": {
        "blocktime": 1231006505,
        "data": "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks",
        "divisible": true,
        "flags": {},
        "issuer": "Satoshi Nakamoto",
        "name": "BTC",
        "propertyid": 0,
        "rdata": null,
        "registered": false,
        "totaltokens": "18316262.00000000",
        "url": "http://www.bitcoin.org"
      },
      "symbol": "BTC",
      "value": "0"
    }
*/
type Addr struct {
	Divisible    bool   `json:"divisible"`
	Error        bool   `json:"error"`
	ID           string `json:"id"`
	Pendingneg   string `json:"pendingneg"`
	Pendingpos   string `json:"pendingpos"`
	Propertyinfo struct {
		Blocktime   int         `json:"blocktime"`
		Data        string      `json:"data"`
		Divisible   bool        `json:"divisible"`
		Flags       struct{}    `json:"flags"`
		Issuer      string      `json:"issuer"`
		Name        string      `json:"name"`
		Propertyid  int         `json:"propertyid"`
		Rdata       interface{} `json:"rdata"`
		Registered  bool        `json:"registered"`
		Totaltokens string      `json:"totaltokens"`
		URL         string      `json:"url"`
	} `json:"propertyinfo"`
	Symbol string `json:"symbol"`
	Value  string `json:"value"`
}

type Balance struct {
	Pendingpos string `json:"pendingpos"`
	Reserved   string `json:"reserved"`
	Divisible  bool   `json:"divisible"`
	Symbol     string `json:"symbol"`
	Value      string `json:"value"`
	Frozen     string `json:"frozen"`
	Pendingneg string `json:"pendingneg"`
	ID         string `json:"id"`
}

type Transaction struct {
	Amount            string `json:"amount"`
	Block             int    `json:"block"`
	Blockhash         string `json:"blockhash"`
	Blocktime         int    `json:"blocktime"`
	Category          string `json:"category"`
	Confirmations     int    `json:"confirmations"`
	Data              string `json:"data"`
	Deadline          int    `json:"deadline"`
	Divisible         bool   `json:"divisible"`
	Earlybonus        int    `json:"earlybonus"`
	Ecosystem         string `json:"ecosystem"`
	Fee               string `json:"fee"`
	Ismine            bool   `json:"ismine"`
	Percenttoissuer   int    `json:"percenttoissuer"`
	Positioninblock   int    `json:"positioninblock"`
	Propertyid        int    `json:"propertyid"`
	Propertyiddesired int    `json:"propertyiddesired"`
	Propertyname      string `json:"propertyname"`
	Propertytype      string `json:"propertytype"`
	Sendingaddress    string `json:"sendingaddress"`
	Subcategory       string `json:"subcategory"`
	Tokensperunit     string `json:"tokensperunit"`
	Txid              string `json:"txid"`
	Type              string `json:"type"`
	TypeInt           int    `json:"type_int"`
	URL               string `json:"url"`
	Valid             bool   `json:"valid"`
	Version           int    `json:"version"`
}

type QueryGethHistory struct {
	Transactions []*Transaction `json:"transactions"`
}
