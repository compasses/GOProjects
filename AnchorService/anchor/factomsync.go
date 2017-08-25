package anchor

import (
	"AnchorService/common"
	"AnchorService/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/go-spew/spew"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var flog = util.FactomLooger

type FactomSync struct {
	factomserver string
	DirBlockMsg  chan common.DirectoryBlockAnchorInfo
}

func NewFactomSync(service *AnchorService) *FactomSync {
	sync := &FactomSync{
		factomserver: service.factomserver,
		DirBlockMsg:  service.DirBlockMsg,
	}

	return sync
}

func (sync *FactomSync) StartSync() {
	// 1. check height and unconfirmed dbblock

	// 2. fetch data of unconfirmed db keyMR and height

	// for mock now
	timeChan := time.NewTicker(time.Second * 10).C

	for {
		select {
		case <-timeChan:
			flog.Info("Got new block info, anchor it...")
			h, err := common.HexToHash("32ce948a6e45cb7e5d098b7c53fe0f60fda14667ac9457bdbafcea04b673918d")
			if err != nil {
				flog.Info("hash error ", err)
				continue
			}
			info := common.DirectoryBlockAnchorInfo{
				KeyMR:    h,
				DBHeight: 556,
			}

			sync.DirBlockMsg <- info
		}
	}

}

func (sync *FactomSync) SyncUp() error {
	// get heights
	heightReq := factom.NewJSON2Request("heights", 0, nil)
	heightResp, err := DoFactomReq(heightReq, sync.factomserver)
	if err != nil {
		flog.Error("call get hegiht error, no need sync up now")
		return fmt.Errorf("error %s", err)
	}
	var result factom.HeightsResponse
	err = json.Unmarshal(heightResp.Result, &result)
	if err != nil {
		flog.Error("Unmarshal error no need sync up now")
		return fmt.Errorf("error %s", err)
	}
	height := result.DirectoryBlockHeight
	// start anchor the top 100
	flog.Info("Start sync up ", " total height", height)
	cfg := util.ReadConfig()
	interval := time.Duration(cfg.Anchor.Interval) * time.Minute
	flog.Info("Sync interval", "interval ", interval)

	timeChan := time.NewTicker(interval).C

	height = height //
	for {
		select {
		case <-timeChan:

			dblock, err := sync.GetDBlockInfoByHeight(height)
			if err != nil {
				flog.Error("Sync up error, get new block error, just check in next round", "err", err)
				continue
			}

			flog.Info("Got dblockanchor info let's anchor it ", "block ", dblock, " height", height)
			sync.DirBlockMsg <- *dblock
			height++
		}
	}

	return nil

}

func (sync *FactomSync) GetDBlockInfoByHeight(height int64) (*common.DirectoryBlockAnchorInfo, error) {
	params := struct {
		Height int64 `json:"height"`
	}{
		Height: height,
	}

	req := factom.NewJSON2Request("dblock-by-height", 0, params)
	resp, err := DoFactomReq(req, sync.factomserver)
	if resp.Error != nil {
		return nil, fmt.Errorf("dblock-by-height error happen %s", resp.Error.Message)
	}

	var dblock = struct {
		Dblock common.DBlockForAnchor
	}{}

	err = json.Unmarshal(resp.Result, &dblock)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error ", err)
	}
	flog.Debug("got dblock ", "block ", spew.Sdump(dblock))

	h, err := common.HexToHash(dblock.Dblock.KeyMR)
	if err != nil {
		return nil, fmt.Errorf("hash error %s", err)
	}

	return &common.DirectoryBlockAnchorInfo{
		KeyMR:    h,
		DBHeight: dblock.Dblock.Header.DBHeight,
	}, nil
}

func DoFactomReq(req *factom.JSON2Request, factomserver string) (*factom.JSON2Response, error) {
	j, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var client *http.Client
	var httpx string

	client = &http.Client{}
	httpx = "http"

	re, err := http.NewRequest("POST",
		fmt.Sprintf("%s://%s/v2", httpx, factomserver),
		bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	re.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(re)
	if err != nil {
		errs := fmt.Sprintf("%s", err)
		if strings.Contains(errs, "\\x15\\x03\\x01\\x00\\x02\\x02\\x16") {
			err = fmt.Errorf("Factomd API connection is encrypted. Please specify -factomdtls=true and -factomdcert=factomdAPIpub.cert (%v)", err.Error())
		}
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("Factomd username/password incorrect.  Edit factomd.conf or\ncall factom-cli with -factomduser=<user> -factomdpassword=<pass>")
	}
	r := factom.NewJSON2Response()
	if err := json.Unmarshal(body, r); err != nil {
		return nil, err
	}

	return r, nil
}
