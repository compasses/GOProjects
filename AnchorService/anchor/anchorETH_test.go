package anchor_test

import (
	"AnchorService/anchor"
	"AnchorService/common"
	"AnchorService/util"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

var (
	AccountAddress    = "0x100c8b406978a413c4305b3AA6074F734feE6C9c"
	AccountPassphrase = "Initial0"
	EthHttpHost       = "localhost:8545"
	GasPrice          = "0x1"
)

func TestAnchorETH_StrToInt(t *testing.T) {
	s1 := "0x15ac8b"

	h, e := strconv.ParseInt(s1, 0, 0)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println("h ", h)
}

func TestMaxTransaction(t *testing.T) {
	Unlock()

	hash, err := common.HexToHash("c6d571c489c76346a45c271fa2ff831036cd7053a1eba37789b3754d20bfb1ea")
	if err != nil {
		fmt.Println("HextoHash error %s", err)
	}

	data, err := anchor.PrependBlockHeight(1222, hash.GetBytes())
	if err != nil {
		fmt.Println(err)
	}

	temp := ""

	for i := 1; i < 10; i++ {
		temp = temp + hex.EncodeToString(data)
	}

	hexs := "0x" + temp

	fmt.Println("got hex string ", hexs)

	txreqJson := util.NewJSON2Request("eth_sendTransaction", 1, []interface{}{
		map[string]interface{}{
			"from":     AccountAddress,
			"to":       AccountAddress,
			"gasPrice": GasPrice,
			"value":    "0x1",
			"data":     hexs,
		},
	})

	fmt.Println("do request ", txreqJson)
	txreq, err := util.EncodeJSON(txreqJson)
	if err != nil {
		fmt.Println(" encode error %s", err)
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("http://%s", EthHttpHost), bytes.NewBuffer(txreq))
	if err != nil {
		fmt.Println("Http New Request error %s", err)
	}

	httpClient := http.DefaultClient
	resp, err := httpClient.Do(httpReq)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Http error %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	r := util.NewJSON2Response()

	if err := json.Unmarshal(body, r); err != nil {
		fmt.Println("Error on http request parse body %s", err)
	}

	fmt.Println("result is ", r)

	strs := string(r.JSONResult())

	fmt.Println("Got transaction ", strs)
}

func Unlock() {
	unlockReqJson := util.NewJSON2Request("personal_unlockAccount", 1, []interface{}{AccountAddress, AccountPassphrase, 3600})
	unlockReq, err := util.EncodeJSON(unlockReqJson)
	if err != nil {
		fmt.Println("Encode error %s", err)
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("http://%s", EthHttpHost), bytes.NewBuffer(unlockReq))
	if err != nil {
		fmt.Println("Http New Request error %s", err)
	}

	httpClient := http.DefaultClient
	resp, err := httpClient.Do(httpReq)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Http error %s", err)
	}

	defer resp.Body.Close()
}
