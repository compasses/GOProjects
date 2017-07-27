package service

import (
	"github.com/compasses/GOProjects/AnchorService/util"
    "github.com/compasses/GOProjects/AnchorService/common"
    "github.com/FactomProject/factom"
)

var log = util.AnchorLogger

type Anchor interface {
	PlaceAnchor(msg common.DirectoryBlockAnchorInfo)
}

type AnchorService struct {
	DoAnchor        Anchor
	DirBlockMsg     chan common.DirectoryBlockAnchorInfo
	AnchorRecordMsg chan interface{}
}

func NewAnchorService(DirBlockMsg, AnchorRecordMsg chan interface{}) *AnchorService {
	cfg := util.ReadConfig()
	service := new(AnchorService)

	if cfg.App.AnchorTo == 0 {
		btc := new(AnchorBTC)
		err := btc.InitRPCClient()
		if err != nil {
			log.Fatal("Error on init RPC :", err)
		}
		service.DoAnchor = btc
		service.DirBlockMsg = DirBlockMsg
		service.AnchorRecordMsg = AnchorRecordMsg
		return service

	} else if cfg.App.AnchorTo == 1 {
		log.Fatal("Not support ETH currently...")
	} else {
		log.Fatal("Not support this kind of anchor")
	}
	return nil
}

func (service *AnchorService) Start()  {
    log.Info("Start Anchor service...")

    for {
       select {
        case anchorMsg := <-service.DirBlockMsg:
           log.Info("Got anchor msg: ", anchorMsg)
           go service.DoAnchor.PlaceAnchor(anchorMsg)
        }
    }
}

func NewEntry(chainid string, external, content []byte) *factom.Entry{
    entry := new(factom.Entry)
    entry.ChainID = chainid
    entry.Content = content
    bta := make([][]byte, 0)
    bta = append(bta, external)
    entry.ExtIDs = bta

    return entry
}

