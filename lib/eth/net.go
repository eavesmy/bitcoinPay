package eth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eavesmy/bitcoinPay/lib"
	"io/ioutil"
	"net/http"
	"strings"
)

const host = "https://api.etherscan.io/api?"
const etherscan_API_Token = "VJTXVRGIKSWYJN3Q3GMJ3YJTJS55IE9S8J"

type etherscan_data struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}
type etherscan_error_data struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type etherscan_error struct {
	Error etherscan_error_data `json:"error"`
}

/*
{"blockNumber":"5856039","timeStamp":"1529994958","hash":"0xcfa582e8720fa3170574d598c6ed0f8ef626988bd1114ba8adb7e8677645e851","nonce":"392305","blockHash":"0x577f13422aac588ae606062e7d22631832bebf9451356cde44037b2d8d93f5bb","from":"0x0d0707963952f2fba59dd06f2b425ace40b492fe","contractAddress":"0xdac17f958d2ee523a2206206994597c13d831ec7","to":"0xffea26119c2cf4d123f2e5ad1c070798ba1352c5","value":"55000000","tokenName":"Tether USD","tokenSymbol":"USDT","tokenDecimal":"6","transactionIndex":"123","gas":"200000","gasPrice":"6000000000","gasUsed":"53465","cumulativeGasUsed":"4191844","input":"deprecated","confirmations":"3630446"}
*/

func GetEthBalance(address string) string {

	data, err := getRequest(map[string]string{
		"module":  "account",
		"action":  "balance",
		"address": address,
		"tag":     "latest",
	})
	if err != nil {
		return "0"
	}
	return data.Result.(string)
}

func GetTokenBalance(addr string, contract string) string {

	data, err := getRequest(map[string]string{
		"module":  "account",
		"action":  "tokenbalance",
        "contractaddress": contract,
		"address": addr,
	})

	if err != nil {
		return "0"
	}

	return data.Result.(string)
}

// Get nonce
func GetTransactionCount(address string) string {

	data, err := getRequest(map[string]string{
		"module":  "proxy",
		"action":  "eth_getTransactionCount",
		"address": address,
	})

	if err != nil {
		return "0x"
	}

	return data.Result.(string)
}

func SendRawTransaction(hex string) string {
	data, err := getRequest(map[string]string{
		"module": "proxy",
		"action": "eth_sendRawTransaction",
		"hex":    hex,
	})
	if err != nil {
		return err.Error()
	}

	return data.Result.(string)
}

func GasPrice() string {
	data, err := getRequest(map[string]string{
		"module": "proxy",
		"action": "eth_gasPrice",
	})
	if err != nil {
		return "0x"
	}

	return data.Result.(string)
}

func GetTransactionByHash(hash string) *lib.Transaction {
	data, err := getRequest(map[string]string{
		"module": "proxy",
		"action": "eth_getTransactionByHash",
		"txhash": hash,
	})
	if err != nil {
		return &lib.Transaction{}
	}

	b, _ := json.Marshal(data.Result)

	tx := &lib.Transaction{}
	json.Unmarshal(b, tx)

	return tx
}

// params: start end sort
func GetTransactions(params map[string]string) []*lib.Transaction {
	data, err := getRequest(map[string]string{
		"module":     "account",
		"action":     "tokentx",
		"startblock": params["start"],
		"endblock":   params["end"],
		"sort":       params["sort"],
		"address":    params["address"],
	})
	if err != nil {
		return []*lib.Transaction{}
	}

	var ret []*lib.Transaction

	rets := data.Result.([]interface{})

	for _, item := range rets {
		b, _ := json.Marshal(item)

		tx := &lib.Transaction{}

		json.Unmarshal(b, tx)

		ret = append(ret, tx)
	}

	return ret
}

// params: [ module action apkKey ...]
func getRequest(params map[string]string) (*etherscan_data, error) {

	params["apikey"] = etherscan_API_Token

	querys := []string{}
	query := ""

	for k, v := range params {
		querys = append(querys, k+"="+v)
	}
	query = strings.Join(querys, "&")

	res, err := http.Get(host + query)

	if err != nil {
		return nil, err
	}

	buffer, _ := ioutil.ReadAll(res.Body)

	data := &etherscan_data{}
	json.Unmarshal(buffer, data)

	if data.Result == "" {
		fmt.Println("error: ", params["module"], params["action"], string(buffer))
		errData := &etherscan_error{}
		json.Unmarshal(buffer, errData)

		return nil, errors.New(errData.Error.Message)
	}

	return data, nil
}
