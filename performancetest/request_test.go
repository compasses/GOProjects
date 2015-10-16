package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestReqHttp(t *testing.T) {
	done := make(chan bool)

	go func() {
		url := "http://cnpvgvb1ep015.pvgl.sap.corp:49192/products/not?eshop_id=25"
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}

		resp.Body.Close()
		log.Print(resp)
		done <- true
	}()
	<-done
}

func TestConfig(t *testing.T) {
	conf := LoadConfig()
	log.Println("config: ", conf)
	httpRequests := BuildHttpRequest(*conf)
	FinalResult := NewStatsRecord(httpRequests, 0, 0)
	log.Println(FinalResult)
	log.Println("get proxy ", os.Getenv("HTTPS_PROXY"))
	proxyUrl, err := url.Parse(os.Getenv("HTTPS_PROXY"))
	log.Printf("Get proxy failed...%v", proxyUrl)

	if err != nil {
		log.Fatalln("Get proxy failed...", err, proxyUrl)
		return
	}
}
