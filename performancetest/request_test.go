package main

import (
	"log"
	"net/http"
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

}
