package main

import (
	"log"
	"net/http"
	"time"
	"github.com/Compasses/Projects-of-GO/apiservice/offline"
)



var GlobalServerStatus int64 = 0
type ServerSwitch map[string]http.Handler

// Implement the ServerHTTP method on our new type
func (sw ServerSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//ChangeGlobalStatus()
	if GlobalServerStatus == 0 { //offline
		sw["offline"].ServeHTTP(w, r)
		log.Println("server run in offline!")
	} else {
		http.Redirect(w, r, "http://10.128.163.72:8080", 302)
	}
}

func ChangeGlobalStatus() {
	if GlobalServerStatus == 0 {
		GlobalServerStatus = 1
	} else {
		GlobalServerStatus = 0
	}
}

func PrintStatus() {
	for {
		log.Println("status", GlobalServerStatus)
		time.Sleep(time.Millisecond * 1000)
	}
}

func main() {

	sw := make(ServerSwitch)

	router := offline.ServerRouter()

	sw["offline"] = router
	//go PrintStatus()

	log.Fatal(http.ListenAndServe(":8080", sw))
}
