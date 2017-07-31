package anchor

import (
	"AnchorService/common"
	"time"
)

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
			log.Info("Got new block info, anchor it...")
			h, err := common.HexToHash("32ce948a6e45cb7e5d098b7c53fe0f60fda14667ac9457bdbafcea04b673918d")
			if err != nil {
				log.Info("hash error ", err)
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
