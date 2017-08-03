package anchor_test

import (
	"AnchorService/anchor"
	"AnchorService/common"
	"encoding/json"
	"fmt"
	"github.com/FactomProject/factom"
	"testing"
)

var factomServer = "localhost:8088"

func TestFactomSync_SyncUp(t *testing.T) {
	heightReq := factom.NewJSON2Request("heights", 0, nil)
	heightResp, err := anchor.DoFactomReq(heightReq, factomServer)
	if err != nil {
		fmt.Println("call get hegiht error, no need sync up now")
		return
	}
	var result factom.HeightsResponse
	err = json.Unmarshal(heightResp.Result, &result)
	if err != nil {
		fmt.Println("Unmarshal error no need sync up now")
		return
	}
	height := result.DirectoryBlockHeight
	// start anchor the top 100
	fmt.Println("Start sync up from height :", height)

	params := struct {
		Height int64 `json:"height"`
	}{
		Height: height,
	}

	req := factom.NewJSON2Request("dblock-by-height", 0, params)
	resp, err := anchor.DoFactomReq(req, factomServer)
	if resp.Error != nil {
		fmt.Println("dblock-by-height error happen ", resp.Error.Message)
		return
	}

	var dblock = struct {
		Dblock common.DBlockForAnchor
	}{} //map[string]interface{}

	err = json.Unmarshal(resp.Result, &dblock)
	if err != nil {
		fmt.Println("Unmarshal error ", err)
		return
	}

	//fmt.Println("got dblock ", spew.Sdump(dblock))
	//block := dblock["dblock"].(map[string]interface{})
	//keyMr := block["keymr"].(string)
	//header := block["header"].(map[string]interface{})
	//dbheight := header["dbheight"].(float64)

	fmt.Println("keymr", dblock.Dblock.KeyMR, " height", dblock.Dblock.Header.DBHeight)

}
