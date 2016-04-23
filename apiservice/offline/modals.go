package offline

import (
	"encoding/json"
	"log"
	"strconv"
)

type TableId int64
type IDSeqs []TableId
type StrSeqs []string

type ATSReq struct {
	SkuIds   IDSeqs
	ChanelId TableId
}

//used to avoid recurision
type idseqs IDSeqs

func (ids *IDSeqs) UnmarshalJSON(b []byte) (err error) {
	log.Println("IDSeqs got bytes: ", string(b))
	var tids *idseqs
	if err = json.Unmarshal(b, &tids); err == nil {
		*ids = IDSeqs(*tids)
		return
	}

	strid := new(StrSeqs)
	if err = json.Unmarshal(b, strid); err == nil {
		for _, val := range *strid {
			if v, err := strconv.ParseInt(val, 10, 64); err == nil {
				*ids = append(*ids, TableId(v))
			}
		}
		return
	} else {
		log.Println("got error: ", err)
	}
	return
}

func (id *TableId) UnmarshalJSON(b []byte) (err error) {
	log.Println("TableId got bytes: ", string(b))
	var tabd int64
	if err = json.Unmarshal(b, &tabd); err == nil {
		*id = TableId(tabd)
		return
	}

	s := ""
	if err = json.Unmarshal(b, &s); err == nil {
		v, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			*id = TableId(v)
		}
	}

	return
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
	CustomerType string  `json:"customerType"`
	Email        string  `json:"email"`
}

type CustomerCreateRsp struct {
	CustomerCode     string  `json:"customerCode"`
	CustomerID       TableId `json:"customerID"`
	ChannelAccountID TableId `json:"channelAccountID"`
	//FailType		 string `json:"failType"`
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

type CartItems struct {
	SkuId              interface{} `json:"skuId"`
	UnitPrice          interface{} `json:"unitPrice"`
	Quantity           interface{} `json:"quantity"`
	TaxAmount          interface{} `json:"taxAmount"`
	DiscountPercentage interface{} `json:"discountPercentage"`
	LineTotal          interface{} `json:"lineTotal"`
	LineTotalAfterDisc interface{} `json:"lineTotalAfterDisc"`
	StandardPrice      interface{} `json:"standardPrice"`
}

type ShoppingCart struct {
	CartTotal          interface{} `json:"cartTotal"`
	DiscountPercentage interface{} `json:"discountPercentage"`
	DiscountSum        interface{} `json:"discountSum"`
	PriceMethod        interface{} `json:"priceMethod"`
	CartItems          []CartItems `json:"cartItems"`
}

type CheckoutCartPlayLoad struct {
	ShippingAddress    interface{}  `json:"shippingAddress"`
	BillingAddress     interface{}  `json:"billingAddress"`
	CustomerId         interface{}  `json:"customerId"`
	ChannelAccountId   interface{}  `json:"channelAccountId"`
	ChannelId          interface{}  `json:"channelId"`
	ShoppingCart       ShoppingCart `json:"shoppingCart"`
	ShippingMethod     interface{}  `json:"shippingMethod"`
	Promotion          interface{}  `json:"promotion"`
	TaxTotal           interface{}  `json:"taxTotal"`
	OrderTotal         interface{}  `json:"orderTotal"`
	DiscountPercentage interface{}  `json:"discountPercentage"`
	DiscountSum        interface{}  `json:"discountSum"`
}

type CheckoutShoppingCart struct {
	ShoppingCart CheckoutCartPlayLoad `json:"shoppingCart"`
}

type CheckoutShoppingCartRsp struct {
	CheckoutCartPlayLoad
	ShippingCosts         interface{} `json:"shippingCosts"`
	EnableExpressDelivery bool        `json:"enableExpressDelivery"`
}

type OrderCreate struct {
	EShopOrder CheckoutCartPlayLoad `json:"eShopOrder"`
}
