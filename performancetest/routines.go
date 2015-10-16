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
	RawURL         string
}

type Config struct {
	Duration    int64
	ThreadNum   []int64
	TestRequest []Requests
}

var FileToSave = "./testresult.txt"

type FinalStats map[string]*RequestStats

func init() {
	f, err := os.OpenFile(FileToSave, os.O_WRONLY|os.O_CREATE, 0x777)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	headline := "Threads" + "\t" + "NumReqs" + "\t" + "TPS" + "\t" + "AvgResp(s)" + "\t" +
		"MaxResp(s)" + "\t" + "MinResp(s)" + "\t" + "ErrNums" + "\t" +
		"ErrCodes" + "\t" + "URL" + "\n"
	f.WriteString(headline)
}

func NewStatsRecord(req RoutineRequest, threads int64, duration int64) FinalStats {
	result := make(FinalStats)
	for _, r := range req {
		result[r.URL] = &RequestStats{}
		result[r.URL].ErrorStatusCode = make(map[int]int)
		result[r.URL].Threads = threads
		result[r.URL].MinResponseTime = 3600
		result[r.URL].Duration = duration
	}
	return result
}

func Collect(req RoutineRequest, threads int64, timeout time.Duration, done chan bool) {
	result := make(chan RequestResult, threads)
	quit := make(chan bool, threads)

	timetick := time.After(timeout * time.Second)
	log.Println("start test ", req, "threads ", threads)

	for i := int64(0); i < threads; i++ {
		go req.StartRoutine(result, quit)
	}

	FinalResult := NewStatsRecord(req, threads, int64(timeout))

	for {
		select {
		case res := <-result:
			FinalResult.StatsRecord(res)
		case <-timetick:
			for i := int64(0); i < threads; i++ {
				quit <- true
			}
			FinalResult.Searilize()
			done <- true
			return
		}
	}
}

func StartCollect(conf Config) {
	done := make(chan bool)
	httpRequests := BuildHttpRequest(conf)

	for _, threads := range conf.ThreadNum {
		go Collect(httpRequests, threads, time.Duration(conf.Duration), done)
		<-done
	}
}

func (result FinalStats) StatsRecord(single RequestResult) {
	stats := result[single.RawURL]

	stats.URL = single.RawURL
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

func (result FinalStats) Searilize() {
	f, err := os.OpenFile(FileToSave, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0x777)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, stats := range result {
		stats.AvgResponseTime = float64(stats.ResponseTime) / float64(stats.NumRequest)
		stats.TPS = float64(stats.NumRequest) / float64(stats.Duration)

		res := fmt.Sprint(stats.Threads) + "\t" + fmt.Sprint(stats.NumRequest) + "\t" +
			fmt.Sprintf("%.3f", stats.TPS) + "\t" +
			fmt.Sprintf("%.3f", stats.AvgResponseTime) + "\t" +
			fmt.Sprintf("%.3f", stats.MaxResponseTime) + "\t" +
			fmt.Sprintf("%.3f", stats.MinResponseTime) + "\t" +
			fmt.Sprint(stats.ErrorNumbers) + "\t" +
			fmt.Sprint(stats.ErrorStatusCode) + "\t" + fmt.Sprint(stats.URL) + "\n"

		f.WriteString(res)
	}

}
