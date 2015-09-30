package offline

type TableId int64

type ATSReq struct {
	SkuIds    []TableId
	ChannelId int
}

type ATSRsp struct {
	SkuId          TableId `json:"skuId"`
	Ats            int64   `json:"ats"`
	AllowBackOrder bool    `json:"allowBackOrder"`
}

type RecommandInfo struct {
	ChannelId  int
	ProductId  TableId `json:",string"`
	CurrencyId int
}

type Customer struct {
	Id           TableId `json:"id,string"`
	CustomerType string
	Email        string
}

type CustomerCreateRsp struct {
	CustomerCode     string  `json:"customerCode"`
	CustomerID       TableId `json:"customerID"`
	ChannelAccountID TableId `json:"channelAccountID"`
}

type CustomerAddress struct {
	Id            TableId     `json:"id"`
	CustomerInfo  Customer    `json:"customer"`
	AddressInfo   interface{} `json:"address"`
	DefaultBillTo bool        `json:"defaultBillTo"`
	DefaultShipTo bool        `json:"defaultShipTo"`
}

type CustomerCreate struct {
	ChannelId    TableId
	Account      string
	Customer     Customer
	CustomerType string
	AccountInfo  CustomerCreateRsp
	Addresses    []CustomerAddress
}
