package btc

/*
   这里封装的 https://api.omniexplorer.info/ api
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	gtype "github.com/eavesmy/golang-lib/type"
	"io/ioutil"
	"net/http"
	"strings"
)

const HOST = "https://api.omniexplorer.info"

// routes
const (
	ADDRDETAIL  = "/v1/address/addr/"
	ADDRBALANCE = "/v2/address/addr/"
	GETHHISTORY = "/v1/properties/gethistory/3"
)

func request(path string, data string) (info []byte) {
	res, err := http.Post(HOST+path, "application/x-www-form-urlencoded", bytes.NewBufferString(data))
	if err != nil {
		fmt.Println(err)
		request(path, data)
		return
	}

	defer res.Body.Close()

	info, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return
}

func GetBalance(addr string, id int) int64 {
	query := "addr=" + addr

	var info map[string][]*Addr
	res := request(ADDRDETAIL, query)

	json.Unmarshal(res, &info)

	str_id := gtype.Int2String(id)

	for _, item := range info["balance"] {
		if item.ID == str_id {
			return gtype.String2Int64(item.Value)
		}
	}

	return 0
}

func GetBalances(addrs []string, id string) map[string]string {

	querys := []string{}
	for _, addr := range addrs {
		querys = append(querys, "addr="+addr)
	}
	query := strings.Join(querys, "&")

	var info map[string]map[string][]*Balance
	res := request(ADDRBALANCE, query)
	json.Unmarshal(res, &info)

	ret := map[string]string{}

	for k, v := range info {
		_v := v["balance"]

		for _, balance := range _v {
			if id == balance.ID {
				ret[k] = balance.Value
			}
		}
	}
	return ret
}

func GetFee() string {

	query := "page=0"
	res := request(GETHHISTORY, query)

	ts := &QueryGethHistory{}
	json.Unmarshal(res, &ts)

	return ts.Transactions[0].Fee
}
