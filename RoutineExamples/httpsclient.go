package main

import (
)
import (
	"net/http"
	"crypto/tls"
	"time"
	"net"
	"io/ioutil"
	"fmt"
	"compress/gzip"
	"net/url"
	"crypto/x509"
	"os"
)

func main() {
	//client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	//if err != nil {
	//	fmt.Println("%v", err)
	//}
	//
	//indexName	:= "stores"
	//routingId	:= 10
	//channelId 	:= 101
	proxyUrl, err := url.Parse("http://proxy.pal.sap.corp:8080")
	cert, err := ioutil.ReadFile("certificate.pem")
	if err != nil {
		fmt.Println("Error occured", err)
		os.Exit(1)
	}
	certp, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config {
			Certificates:       []tls.Certificate{certp},
			InsecureSkipVerify: true,
			RootCAs: certPool},
		Proxy: http.ProxyURL(proxyUrl),
		DisableCompression: true,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	newRq, err := http.NewRequest("GET", "https://sigma.sapanywhere.io", nil)
	if err != nil {
		fmt.Println("new request error ", err)
	}


	resp, err := client.Do(newRq)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		fmt.Println("GET response error ", err)
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Body, err = gzip.NewReader(resp.Body)
		if err != nil {
			panic(err)
		}
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil read err ", err)
	}

	fmt.Printf("Get response %s\n", res)
}