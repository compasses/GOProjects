package main

import (
	"fmt"
	"time"
)

func TimeOutMsg() string {
	time.Sleep(time.Second)
	return "msg time out"
}

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			timeout := false
			select {
			case <-time.After(time.Second):
				fmt.Println("time out msg 1")
				timeout = true
			case c1 <- TimeOutMsg():
				time.Sleep(time.Second * 3)
				fmt.Println("sending msg 1")
			}

			if timeout {
				break
			}
		}
	}()

	go func() {
		for {
			timeout := false
			select {
			case <-time.After(time.Second):
				fmt.Println("time out msg 2")
				timeout = true
			case c2 <- TimeOutMsg():
				time.Sleep(time.Second * 3)
				fmt.Println("sending msg 2")
			}
			if timeout {
				break
			}
		}

	}()

	go func() {
		for {
			timeout := false
			select {
			case msg1 := <-c1:
				fmt.Println("msg1", msg1)
			case msg2 := <-c2:
				fmt.Println("msg2", msg2)
			case <-time.After(time.Millisecond * 5000):
				fmt.Println("timeout")
				timeout = true
				//			default:
				//				//fmt.Println("nothing ready")
			}
			if timeout {
				break
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}
