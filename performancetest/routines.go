package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type RequestStats struct {
	URL             string
	Threads         int64
	Duration        int64
	TPS             float64
	NumRequest      int64
	ResponseTime    float64
	MaxResponseTime float64
	MinResponseTime float64
	AvgResponseTime float64
	ErrorNumbers    int64
	ErrorStatusCode map[int]int
}

type Requests struct {
	URL    string
	Method string
	Body   string
	Header map[string]string
}

type RequestResult struct {
	SingleRespTime float64
	StatusCode     int
}

type Config struct {
	Duration    int64
	ThreadNum   []int64
	TestRequest []Requests
}

var fileToSave = "./testresult.txt"
var finalResult map[string][]RequestStats

func init() {
	f, err := os.OpenFile(fileToSave, os.O_WRONLY|os.O_CREATE, 0x777)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	headline := "Threads" + "\t" + "NumReqs" + "\t" + "TPS" + "\t" + "AvgResp(s)" + "\t" +
		"MaxResp(s)" + "\t" + "MinResp(s)" + "\t" + "ErrNums" + "\t" +
		"ErrCodes" + "\t" + "URL" + "\n"
	f.WriteString(headline)
	finalResult = make(map[string][]RequestStats)
}

func Collect(req Requests, threads int64, timeout time.Duration, done chan bool) {
	result := make(chan RequestResult, threads)
	quit := make(chan bool, threads)

	client := NewClient(req)
	timetick := time.After(timeout * time.Second)

	log.Println("start test ", client, "threads ", threads)

	for i := int64(0); i < threads; i++ {
		go client.startRoutine(result, quit)
	}

	var stats RequestStats
	stats.URL = req.URL
	stats.Threads = threads
	stats.MinResponseTime = float64(timeout)
	stats.Duration = int64(timeout)
	stats.ErrorStatusCode = make(map[int]int)

	for {
		select {
		case res := <-result:
			stats.StatsRecord(res)
		case <-timetick:
			for i := int64(0); i < threads; i++ {
				quit <- true
			}
			finalResult[req.URL] = append(finalResult[req.URL], stats)
			done <- true
			return
		}
	}
}

func StartCollect(conf Config) {
	done := make(chan bool)

	for _, threads := range conf.ThreadNum {
		for _, req := range conf.TestRequest {
			go Collect(req, threads, time.Duration(conf.Duration), done)
			<-done
		}
	}

	log.Println("Start to write results...")
	for url, _ := range finalResult {
		for _, stat := range finalResult[url] {
			stat.Searilize()
		}
	}
}

func (stats *RequestStats) StatsRecord(single RequestResult) {
	stats.NumRequest++
	stats.ResponseTime += single.SingleRespTime

	if single.SingleRespTime > stats.MaxResponseTime {
		stats.MaxResponseTime = single.SingleRespTime
	}

	if single.SingleRespTime < stats.MinResponseTime {
		stats.MinResponseTime = single.SingleRespTime
	}

	//sum the error code information
	if single.StatusCode != 200 {
		stats.ErrorNumbers++
		stats.ErrorStatusCode[single.StatusCode]++
	}
}

func (stats *RequestStats) Searilize() {
	stats.AvgResponseTime = float64(stats.ResponseTime) / float64(stats.NumRequest)
	stats.TPS = float64(stats.NumRequest) / float64(stats.Duration)

	f, err := os.OpenFile(fileToSave, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0x777)
	if err != nil {
		log.Printf(" stats result %v \n", stats)
		panic(err)
	}

	defer f.Close()
	res := fmt.Sprint(stats.Threads) + "\t" + fmt.Sprint(stats.NumRequest) + "\t" +
		fmt.Sprintf("%.3f", stats.TPS) + "\t" +
		fmt.Sprintf("%.3f", stats.AvgResponseTime) + "\t" +
		fmt.Sprintf("%.3f", stats.MaxResponseTime) + "\t" +
		fmt.Sprintf("%.3f", stats.MinResponseTime) + "\t" +
		fmt.Sprint(stats.ErrorNumbers) + "\t" +
		fmt.Sprint(stats.ErrorStatusCode) + "\t" + fmt.Sprint(stats.URL) + "\n"

	f.WriteString(res)
}
