package online

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
	"strings"
  	"compress/gzip"
)

type ProxyRoute struct {
	client *http.Client
	url    string
	GrabIF string
}

func NewProxyHandler(newurl , grabIF string) *ProxyRoute {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	
	if grabIF != "" {
		go func() {
			for {
				// Wait for 10s.
				time.Sleep(10 * time.Second)
				if (FailNum+SuccNum) > 0 {
					log.Printf("\n\tIF: %s SuccNum:%d FailNum:%d FailureRate:%f\n\n", grabIF, SuccNum,FailNum,float32((FailNum))/float32((FailNum+SuccNum)))
				}
			}
		}()
	}

	return &ProxyRoute{
		client: &http.Client{Transport: tr},
		url:    newurl,
		GrabIF:	grabIF}
}

func (proxy *ProxyRoute) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	newbody := make([]byte, req.ContentLength)
	req.Body.Read(newbody)
	NeedLog := strings.Contains(req.RequestURI, proxy.GrabIF)
	
	newRq, err := http.NewRequest(req.Method, proxy.url+req.RequestURI, ioutil.NopCloser(bytes.NewReader(newbody)))
	if err != nil {
		log.Println("new request error ", err)
	}
	newRq.Header = req.Header

	LogOutPut(NeedLog, "New Request: ")
	RequstFormat(NeedLog, newRq, string(newbody))
	now := time.Now()
	resp, err := proxy.client.Do(newRq)
	if resp != nil {
		defer resp.Body.Close()
	}
	
	LogOutPut(NeedLog, "Time used: ", time.Since(now))
	if err != nil {
		log.Println("get error ", err)
	} else {
		if resp.Header.Get("Content-Encoding") == "gzip" {
		  	resp.Body, err = gzip.NewReader(resp.Body)
            if err != nil {
                    panic(err)
            }
		}
		
		res, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("ioutil read err ", err)
		}
		
		if NeedLog {
			if resp.StatusCode != 200 {
				FailNum++
			} else {
				SuccNum++
			}
		}
		
		LogOutPut(NeedLog, "Get response : ")
		ResponseFormat(NeedLog, resp, string(res))
		for key, _ := range resp.Header {
			w.Header().Set(key, strings.Join(resp.Header[key], ";"))
		}
		
		w.Write(res)
	}
}
