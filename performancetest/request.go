package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RequestClient struct {
	client *http.Client
	req    Requests
}

//get http request client through url, support https
func NewClient(req Requests) *RequestClient {
	u, err := url.Parse(req.URL)
	if err != nil {
		log.Fatal(err)
	}

	if strings.ToLower(u.Scheme) == "https" {
		tr := &http.Transport{}

		return &RequestClient{
			client: &http.Client{Transport: tr},
			req:    req,
		}
	}

	return &RequestClient{
		client: http.DefaultClient,
		req:    req,
	}
}

func buildHeader(req Requests, header *http.Header) {
	header.Set("Connection", "close")
	if strings.ToLower(req.Method) == "post" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	for k, v := range req.Header {
		header.Set(k, v)
	}
}

func buildBody(req Requests) io.Reader {
	if strings.ToLower(req.Method) == "get" || len(req.Body) <= 0 {
		return nil
	}

	newbody := []byte(req.Body)
	return ioutil.NopCloser(bytes.NewReader(newbody))
}

func (httpclient *RequestClient) startRoutine(result chan<- RequestResult, quit chan bool) {

	for {
		request, err := http.NewRequest(httpclient.req.Method, httpclient.req.URL, buildBody(httpclient.req))
		if err != nil {
			log.Fatal(err)
			return
		}

		buildHeader(httpclient.req, &request.Header)
		now := time.Now()
		resp, err := httpclient.client.Do(request)
		elapsed := time.Since(now)

		if err != nil {
			log.Fatalln(err)
		}

		resp.Body.Close()

		select {
		case result <- RequestResult{elapsed.Seconds(), resp.StatusCode}:
		case <-quit:
			return
		}
	}
}
