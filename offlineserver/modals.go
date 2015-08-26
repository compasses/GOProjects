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

type RecommandInfo struct {
	ChannelId  int
	ProductId  int `json:",string"`
	CurrencyId int
}

type Customer struct {
	CustomerType string
	Email        string
}

type CustomerCreateRsp struct {
	CustomerCode     string `json:"customerCode"`
	CustomerID       int64  `json:"customerID"`
	ChannelAccountID int64  `json:"channelAccountID"`
}

type CustomerAddress struct {
	Id            int64
	AddressStr    string
	defaultBillTo bool
	defaultShipTo bool
}

type CustomerCreate struct {
	ChannelId    int64
	Account      string
	Customer     Customer
	CustomerType string
	AccountInfo  CustomerCreateRsp
	Addresses	 []CustomerAddress
}
