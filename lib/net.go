package lib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

const etherscan_API_Token = ""

type etherscan_data struct {
	status  string
	message string
	result  string
}

func GetEthBalance(address string) int64 {
	res, err := http.Get("https://api.etherscan.io/api?module=account&action=balance&address=" + address + "&tag=latest&apikey=" + etherscan_API_Token)

	if err != nil {
		return -1
	}

	buffer, _ := ioutil.ReadAll(res.Body)

	data := &etherscan_data{}

	json.Unmarshal(buffer, data)

	balance, _ := strconv.ParseInt(data.result, 10, 64)

	return balance
}
