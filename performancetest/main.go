package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type RequestStats struct {
	NumRequest      int64
	ResponseTime    int64
	MaxResponseTime int64
	MinResponseTime int64
	ErrorNumbers    int64
	ErrorStatusCode []int64
}

type Requests struct {
	URL string
	Method string
	Body string
	Header map[string]string
}

type Config struct {
	Duration  int64
	ThreadNum []int64
	TestRequest      []Requests
}

func main() {
	//get configuration
	var conf Config
	file, err := os.Open("./config.json")
	if err != nil {
		log.Println("read file failed...", err)
		log.Println("Just run in offline mode")
	} else {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("read file failed...", err)
			log.Println("Just run in offline mode")
		} else {
			err := json.Unmarshal([]byte(string(data)), &conf)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("get configuration:", conf)
		}
	}
	
}
