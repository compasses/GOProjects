package main

import (
	"encoding/json"
	"github.com/Compasses/Projects-of-GO/apiservice/offline"
	"github.com/Compasses/Projects-of-GO/apiservice/online"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	RunMode      int
	RemoteServer string
	ListenOn     string
	LogFile      string
}

var GlobalServerStatus int64 = 0
var localServer string = "localhost:8080"
var GlobalConfig config

func GetConfiguration() (conf config, useDefault bool) {
	//get configuration
	file, err := os.Open("./config.json")
	if err != nil {
		log.Println("read file failed...", err)
		log.Println("Just run in offline mode")
		useDefault = true
	} else {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("read file failed...", err)
			log.Println("Just run in offline mode")
			useDefault = true
		} else {
			json.Unmarshal([]byte(string(data)), &conf)
			log.Println("get configuration:", conf)
		}
	}
	GlobalConfig = conf
	return
}

func init() {

}

func RunDefaultServer(local string, handler http.Handler) {
	log.Println("Listen ON: ", local)
	//log.Fatal(http.ListenAndServe(local, handler))
	log.Fatal(http.ListenAndServeTLS(local, "cert.pem", "key.pem", handler))
}

func main() {
	conf, useDefault := GetConfiguration()

	f, err := os.OpenFile(GlobalConfig.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opening file: %v", err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.SetOutput(f)

	log.Println("Begin API LOG------------------------")

	go func(f *os.File) {
		for {
			f.Sync()
			time.Sleep(time.Second)
		}
	}(f)

	if useDefault {
		router := offline.ServerRouter()
		RunDefaultServer(localServer, router)
	} else {
		if conf.RunMode == 0 {
			log.Println("API Run in offline mode...")
			router := offline.ServerRouter()
			RunDefaultServer(conf.ListenOn, router)
		} else {
			log.Println("API Run in online mode...")
			proxy := online.NewProxyHandler(conf.RemoteServer)
			RunDefaultServer(conf.ListenOn, proxy)
		}
	}

}
