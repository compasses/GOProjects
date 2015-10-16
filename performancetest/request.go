package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type RequestClient struct {
	Client *http.Client
	Requests
	httpReq *http.Request
}

type RoutineRequest []*RequestClient

func BuildHttpRequest(conf Config) (requests RoutineRequest) {
	for _, req := range conf.TestRequest {
		requests = append(requests, NewClient(req))
	}
	return
}

func BuildHeader(req Requests, header *http.Header) {
	header.Set("Connection", "close")
	if strings.ToLower(req.Method) == "post" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	for k, v := range req.Header {
		header.Set(k, v)
	}
}

func BuildBody(req RequestClient) io.Reader {
	if strings.ToLower(req.Method) == "get" || len(req.Body) <= 0 {
		return nil
	}

	newbody := []byte(req.Body)
	return ioutil.NopCloser(bytes.NewReader(newbody))
}

func HandleProxy(isHttps bool) func(*http.Request) (*url.URL, error) {
	var envHttpProxy string

	if isHttps {
		envHttpProxy = "HTTPS_PROXY"
	} else {
		envHttpProxy = "HTTP_PROXY"
	}

	proxyUrl, _ := url.Parse(os.Getenv(envHttpProxy))
	if proxyUrl == nil || proxyUrl.String() == "" {
		log.Println("Get proxy failed...", "not use the proxy...")
		return nil
	}

	return func(*http.Request) (*url.URL, error) {
		return proxyUrl, nil
	}
}

//get http request client through url, support https
func NewClient(req Requests) *RequestClient {
	u, err := url.Parse(req.URL)
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest(req.Method, req.URL, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	BuildHeader(req, &request.Header)

	if strings.ToLower(u.Scheme) == "https" {
		tr := &http.Transport{
			TLSClientConfig:    &tls.Config{},
			DisableCompression: true,
			Proxy:              HandleProxy(true),
		}

		return &RequestClient{
			Client:   &http.Client{Transport: tr},
			Requests: req,
			httpReq:  request,
		}
	}

	tr := &http.Transport{
		Proxy: HandleProxy(false),
	}

	return &RequestClient{
		Client:   &http.Client{Transport: tr},
		Requests: req,
		httpReq:  request,
	}
}

func (httpclient RoutineRequest) StartRoutine(result chan<- RequestResult, quit chan bool) {

	for {
		for _, request := range httpclient {
			req := *request.httpReq //http.Request not thread safe, need copy it.

			now := time.Now()
			resp, err := request.Client.Do(&req)
			elapsed := time.Since(now)

			if err != nil {
				log.Println(err, elapsed)
				if resp != nil {
					log.Println("body ", resp.Status)
				}
				continue
			}

			resp.Body.Close()
			select {
			case result <- RequestResult{elapsed.Seconds(), resp.StatusCode, request.URL}:
			case <-quit:
				return
			}
		}
	}
}
