package main

type ATSReq struct {
	SkuIds    []int
	ChannelId int
}

type ATSRsp struct {
	SkuId          int  `json:"skuId"`
	Ats            int  `json:"ats"`
	AllowBackOrder bool `json:"allowBackOrder"`
}
