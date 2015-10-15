//refter to: http://matt.aimonetti.net/posts/2012/11/27/real-life-concurrency-in-go/

package main

import (
	"fmt"
	"net/http"
	"time"
)

var urls = []string{
	"http://www.baidu.com/",
	"http://www.google.com/",
	"http://www.jd.com/",
}

type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}

func asyncHttpGets(urls []string) (responses []*HttpResponse) {
	ch := make(chan *HttpResponse)

	for _, url := range urls {
		go func(url string) {
			fmt.Println("fetching url ", url)
			resp, err := http.Get(url)
			ch <- &HttpResponse{url, resp, err}
		}(url)
	}

	for {
		select {
		case r := <-ch:
			fmt.Printf("%s was fetched \n", r.url)
			if r.err != nil {
				fmt.Printf("error returned ", r.err)
				return
			}

			responses = append(responses, r)
			if len(responses) == len(urls) {
				return
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
	return
}

func main() {
	results := asyncHttpGets(urls)
	for _, result := range results {
		fmt.Printf("%s status: %s\n", result.url,
			result.response.Status)
	}
}
