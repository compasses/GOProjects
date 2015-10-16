//refer: https://golang.org/ref/mem

package main

import (
	"sync"
	"time"
)

var c = make(chan int, 1)
var a string

func f() {
	a = "hello, world"
	<-c
}

func test() {
	go f()
	c <- 0
	println(a)
}

type worFunc func()

func workwork() {
	time.Sleep(time.Second)
	print(" work work ...")
}

func limitWork() {
	var limit = make(chan int, 3)
	work := make(map[int]worFunc)
	work[0] = worFunc(workwork)
	work[1] = worFunc(workwork)
	work[2] = worFunc(workwork)

	for _, w := range work {
		go func(w worFunc) {
			limit <- 1
			w()
			<-limit
		}(w)
	}

	select {}
}

func flock() {
	a = "hello, world"
	l.Unlock()
}

var l sync.Mutex

func useLock() {
	l.Lock()
	go flock()
	l.Lock()
	print(a)
}

var once sync.Once

func setup() {
	a = "hello, world"
}

func doprint() {
	once.Do(setup)
	print(a)
}

func twoprint() {
	go doprint()
	go doprint()
}

func main() {
	limitWork()
}
