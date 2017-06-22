package main

import "strconv"

func main() {
	go startRPC()
	N := 1000
	mapChan := make(chan int, N)

	for i := 1; i < N; i++ {
		go func(i int) {
			call("localhost", "Worker.Work", strconv.Itoa(i), new(string))
			mapChan <- i
		}(i)
	}
	for i := 0; i < N; i++ {
		<-mapChan
	}
}
