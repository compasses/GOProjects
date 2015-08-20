package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	_ "strconv"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	r.ParseForm() //解析参数，默认是不会解析的

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
