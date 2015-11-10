package offline

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"strings"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	r.ParseForm() //解析参数，默认是不会解析的

	log.Println(net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]))
	dec := json.NewDecoder(r.Body)
	var result interface{}
	dec.Decode(&result)
	log.Println("Req:", result)
}

func PlaceOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var result OrderCreate
	dec.Decode(&result)

	newOrder := RepoCreateOrder(result)

	log.Println("PlaceOrder Rsp: ", newOrder)

	if err := json.NewEncoder(w).Encode(newOrder); err != nil {
		panic(err)
	}
}

func GetSalesOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	dec := json.NewDecoder(r.Body)
	var req interface{}
	dec.Decode(&req)
	Req := req.(map[string]interface{})
	Id := TableId(ToInt64FromString(Req["orderId"].(string)))

	log.Println("GetSalesOrder", req)
	salesOrder := RepoGetSalesOrder(Id, Req["channelAccountId"].(string))

	if err := json.NewEncoder(w).Encode(salesOrder); err != nil {
		panic(err)
	}
}

func GetSalesOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	dec := json.NewDecoder(r.Body)
	var req interface{}
	dec.Decode(&req)
	Req := req.(map[string]interface{})

	log.Println("GetSalesOrders", req)
	salesOrders := RepoGetSalesOrders(Req["channelAccountId"].(string))

	if err := json.NewEncoder(w).Encode(salesOrders); err != nil {
		panic(err)
	}
	log.Println("GetSalesOrders Rsp", salesOrders)

}

func Checkout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var result CheckoutShoppingCart
	err := dec.Decode(&result)
	if err != nil {
		HandleError(err)
	}

	resp := RepoCheckoutShoppingCart(result.ShoppingCart)
	log.Println("CheckoutShoppingCart resp: ", resp)
	
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}

}

func ATS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var checkInfo ATSReq
	err := dec.Decode(&checkInfo)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("ATS Req %+v\n", checkInfo)
	rsp := RepoCreateATSRsp(&checkInfo)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		panic(err)
	}
}

func RecommandationProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var id RecommandInfo
	err := dec.Decode(&id)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("RecommandProducts Req%+v\n", id)

	RecommandIds := RepoCreateRecommandationProducts(id.ProductId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(RecommandIds); err != nil {
		panic(err)
	}
	log.Printf("RecommandProducts Rsp%+v\n", RecommandIds)
}

func GetCustomer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	var channelAccountId interface{}
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&channelAccountId)
	id := channelAccountId.(map[string]interface{})

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Account := RepoGetCustomer(GetIdFromStr(id["channelAccountId"].(string)))

	if err := json.NewEncoder(w).Encode(Account); err != nil {
		panic(err)
	}
	log.Printf("Return Customer Rsp %+v\n", Account)
}

func CreateCustomer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	var customer CustomerCreate
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&customer)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Customer Create Req%+v\n", customer)

	Account := RepoCreateAccount(customer)

	if err := json.NewEncoder(w).Encode(Account); err != nil {
		panic(err)
	}
	log.Printf("Customer Create Rsp%+v\n", Account)
}

func CustomerAddressNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)

	var addInfo CustomerAddress
	err := dec.Decode(&addInfo)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Rs := RepoCreateAddress(&addInfo)

	if err = json.NewEncoder(w).Encode(Rs); err != nil {
		panic(err)
	}

	log.Printf("Create address info %+v\n", addInfo)
	log.Printf("Result %+v\n", Rs)

}

func CustomerAddressUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	addressId := GetIdFromStr(ps.ByName("id"))

	dec := json.NewDecoder(r.Body)

	var addInfo CustomerAddress
	err := dec.Decode(&addInfo)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Rs := RepoUpdateAddress(addressId, &addInfo)
	log.Printf("Update address info %+v\n", Rs)

	if err = json.NewEncoder(w).Encode(Rs); err != nil {
		panic(err)
	}
}

func GetCustomerAddress(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	customerId := GetIdFromStr(r.Form["$filter"][0])
	Rs := RepoGetCustomerAddress(customerId)

	if err := json.NewEncoder(w).Encode(Rs); err != nil {
		panic(err)
	}
}
func MiscCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)

	var checkParam map[string]interface{}
	err := dec.Decode(&checkParam)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("misc check parames ", checkParam)
	Rs := RetrieveByMapLevel(checkParam, []string{"miscParam", "lines"})
	lines := Rs.([]interface{})

	resp := make(map[string][]interface{})
	for _, val := range lines {
		valm := val.(map[string]interface{})
		resp["lineResult"] = append(resp["lineResult"], map[string]interface{}{
			"onChannel":      "true",
			"ats":            10,
			"allowBackOrder": "true",
			"skuId":          valm["skuId"],
			"valid":          "true",
		})
	}

	log.Println("resp ", resp)

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}
