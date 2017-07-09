package BechMark

import "sync"

var mutex = sync.Mutex{}
var tchan = make(chan bool, 1)

func ForBenchMutext() {
    mutex.Lock()
    mutex.Unlock()
}

func ForBenchChan() {
    tchan <- true
    <- tchan
}



func main() {
    inChan := make(chan int)

    inChan <- 1

    println(<-inChan)
}
