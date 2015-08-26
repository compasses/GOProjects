package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	_ "strconv"
	"strings"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	r.ParseForm() //解析参数，默认是不会解析的

	fmt.Println(net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]))

	all, _ := ioutil.ReadAll(r.Body)
	var result interface{}
	json.Unmarshal(all, &result)
	fmt.Println("Req:", result)
}

func ATS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()

	body, _ := ioutil.ReadAll(r.Body)

	var checkInfo ATSReq
	err := json.Unmarshal(body, &checkInfo)

	if err != nil {
		fmt.Println("err :", err)
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
	body, _ := ioutil.ReadAll(r.Body)

	var id RecommandInfo
	err := json.Unmarshal(body, &id)

	if err != nil {
		fmt.Println("err :", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("RecommandProducts Req%+v\n", id)

	key := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])

	RecommandIds := RepoCreateRecommandationProducts(key, id.ProductId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(RecommandIds); err != nil {
		panic(err)
	}
	log.Printf("RecommandProducts Rsp%+v\n", RecommandIds)
}

func CreateCustomer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	body, _ := ioutil.ReadAll(r.Body)
	var customer CustomerCreate
	err := json.Unmarshal(body, &customer)

	if err != nil {
		fmt.Println("err :", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Customer Create Req%+v\n", customer)

	key := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])

	Account := RepoCreateAccount(key, customer)

	if err := json.NewEncoder(w).Encode(Account); err != nil {
		panic(err)
	}
	log.Printf("Customer Create Rsp%+v\n", Account)
}
