package online

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type ProxyRoute struct {
	client *http.Client
	url    string
}

func NewProxyHandler(newurl string) *ProxyRoute {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &ProxyRoute{
		client: &http.Client{Transport: tr},
		url:    newurl}
}

func (proxy *ProxyRoute) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	newbody := make([]byte, req.ContentLength)
	req.Body.Read(newbody)

	newRq, err := http.NewRequest(req.Method, proxy.url+req.RequestURI, ioutil.NopCloser(bytes.NewReader(newbody)))
	if err != nil {
		log.Println("new request error ", err)
	}
	newRq.Header = req.Header

	log.Println("New Request: ")
	RequstFormat(newRq, string(newbody))
	now := time.Now()
	resp, err := proxy.client.Do(newRq)
	defer resp.Body.Close()
	log.Println("Time used: ", time.Since(now))
	if err != nil {
		log.Println("get error ", err)
	} else {
		res, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("ioutil read err ", err)
		}
		log.Println("Get response : ")
		ResponseFormat(resp, string(res))
		w.Write(res)
	}
}
