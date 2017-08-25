package anchor

import (
	"AnchorService/common"
	"AnchorService/util"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FactomProject/factomd/anchor"
	"github.com/FactomProject/go-spew/spew"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type AnchorETH struct {
	accountAddress string
	accountPass    string
	ethHost        string
	gasprice       string
	service        *AnchorService
	count          int
}

func NewAnchorETH() *AnchorETH {
	eth := &AnchorETH{}
	return eth
}

func (eth *AnchorETH) InitEthClient() error {
	log.Info("Init anchor eth")
	cfg := util.ReadConfig()
	eth.accountAddress = cfg.Eth.AccountAddress
	eth.accountPass = cfg.Eth.AccountPassphrase
	eth.ethHost = cfg.Eth.EthHttpHost
	eth.gasprice = cfg.Eth.GasPrice

	return nil
}

func (anchorETH *AnchorETH) PlaceAnchor(msg common.DirectoryBlockAnchorInfo) {
	anchorRec := new(anchor.AnchorRecord)
	eth := new(anchor.EthereumStruct)

	anchorRec.Ethereum = eth
	anchorRec.KeyMR = msg.KeyMR.String()
	anchorRec.DBHeight = msg.DBHeight
	anchorRec.AnchorRecordVer = 1
	anchorETH.count++

	if err := anchorETH.doTransaction(anchorRec); err != nil {
		anchorETH.service.AnchorFail <- false
	}
}

func (anchorETH *AnchorETH) doTransaction(record *anchor.AnchorRecord) error {
	// unlock account
	err := anchorETH.unlockAccount()
	if err != nil {
		log.Error("Cannot unlock account, cannot anchor now...")
		return errors.New("Error")
	}

	// do transaction with data
	txHashStr, err := anchorETH.sendTransaction(record)
	if err != nil {
		log.Error("Send transaction error ", err)
		return errors.New("Error")
	}

	log.Info("Got transaction", "hash ", *txHashStr)
	timeChan := time.NewTicker(time.Minute).C
	totalTry := 60
	// wait confirm to save anchor into factom
ForLoop:
	for {
		select {
		case <-timeChan:
			totalTry--
			if totalTry <= 0 {
				break ForLoop
			}

			receipt, err := anchorETH.getTransactionReceipt(*txHashStr)
			if err != nil {
				log.Info("error happen ", "got error", err, "total try left ", totalTry)
				continue
			}
			log.Debug("Got receipt ", "info", spew.Sdump(receipt))

			anchorETH.saveAnchor(receipt, record)
			break ForLoop
		}
	}

	if totalTry <= 0 {
		return errors.New("Error")
	}
	return nil
}

func (eth *AnchorETH) saveAnchor(receipt *common.EthTxReceipt, record *anchor.AnchorRecord) {
	record.Ethereum.TXID = receipt.TransactionHash
	record.Ethereum.BlockHash = receipt.BlockHash
	record.Ethereum.Address = receipt.From
	record.Ethereum.BlockHeight, _ = strconv.ParseInt(receipt.BlockNumber, 0, 0)
	record.Ethereum.Offset, _ = strconv.ParseInt(receipt.TransactionIndex, 0, 0)

	eth.service.submitEntryToAnchorChain(record)
}

func (eth *AnchorETH) getTransactionReceipt(txHashStr string) (*common.EthTxReceipt, error) {
	log.Info("got receipt hashstr ", "str", txHashStr)
	receiptJson := util.NewJSON2Request("eth_getTransactionReceipt", 1, []interface{}{txHashStr})

	receiptReq, err := util.EncodeJSON(receiptJson)
	if err != nil {
		return nil, err
	}

	log.Debug("receipt", "Receipt get ", spew.Sdump(receiptJson))

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("http://%s", eth.ethHost), bytes.NewBuffer(receiptReq))
	if err != nil {
		return nil, fmt.Errorf("Http New Request error %s", err)
	}

	httpClient := http.DefaultClient
	resp, err := httpClient.Do(httpReq)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("Http error %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	r := util.NewJSON2Response()

	if err := json.Unmarshal(body, r); err != nil {
		return nil, fmt.Errorf("Error on http request parse body %s", err)
	}

	log.Debug("Receipt get result ", "result", spew.Sdump(r))

	if r.Error != nil {
		return nil, fmt.Errorf("Receipt call error...%s", r.Error.Message)
	}
	if string(r.JSONResult()) == "null" {
		return nil, fmt.Errorf("Receipt not generate now, retry later...")
	}

	receipt := common.EthTxReceipt{}

	if err = json.Unmarshal(r.JSONResult(), &receipt); err != nil {
		return nil, fmt.Errorf("Error on unmarshal receipt %s", err)
	}

	return &receipt, nil
}

func (eth *AnchorETH) sendTransaction(record *anchor.AnchorRecord) (*string, error) {
	hash, err := common.HexToHash(record.KeyMR)
	if err != nil {
		return nil, fmt.Errorf("HextoHash error %s", err)
	}

	data, err := PrependBlockHeight(record.DBHeight, hash.GetBytes())
	if err != nil {
		return nil, err
	}
	hexs := "0x" + hex.EncodeToString(data)
	log.Info("got hex string ", "info", hexs)

	txreqJson := util.NewJSON2Request("eth_sendTransaction", 1, []interface{}{
		map[string]interface{}{
			"from":     eth.accountAddress,
			"to":       eth.accountAddress,
			"gasPrice": eth.gasprice,
			"value":    "0x1",
			"data":     hexs,
		},
	})

	txreq, err := util.EncodeJSON(txreqJson)
	if err != nil {
		return nil, fmt.Errorf(" encode error %s", err)
	}

	log.Debug("Send transaction dump ", "dump", spew.Sdump(txreqJson))
	httpReq, err := http.NewRequest("POST", fmt.Sprintf("http://%s", eth.ethHost), bytes.NewBuffer(txreq))
	if err != nil {
		return nil, fmt.Errorf("Http New Request error %s", err)
	}

	httpClient := http.DefaultClient
	resp, err := httpClient.Do(httpReq)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("Http error %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	r := util.NewJSON2Response()

	if err := json.Unmarshal(body, r); err != nil {
		return nil, fmt.Errorf("Error on http request parse body %s", err)
	}

	log.Debug("Got send transaction result ", "result", spew.Sdump(r))

	var result string

	if err := json.Unmarshal(r.JSONResult(), &result); err != nil {
		return nil, fmt.Errorf("Error on parse json result %s", err)
	}

	return &result, nil
}

func (eth *AnchorETH) unlockAccount() error {
	unlockReqJson := util.NewJSON2Request("personal_unlockAccount", 1, []interface{}{eth.accountAddress, eth.accountPass, 3600})
	unlockReq, err := util.EncodeJSON(unlockReqJson)
	if err != nil {
		return fmt.Errorf("Encode error %s", err)
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("http://%s", eth.ethHost), bytes.NewBuffer(unlockReq))
	if err != nil {
		return fmt.Errorf("Http New Request error %s", err)
	}

	httpClient := http.DefaultClient
	resp, err := httpClient.Do(httpReq)
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Http error %s", err)
	}

	defer resp.Body.Close()
	return nil
}
