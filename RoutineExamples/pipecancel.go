//refer to: http://blog.golang.org/pipelines
package main

import (
	"fmt"
	"sync"
)

//first stage generate numbers
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

//second stage output the square number
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

//merge function converts a list of channels to a single channel by starting a goroutine for each
//inbound channel that copies the values to the sole outbound channel. once all the output goroutine have been started
//merge starts one more goroutine to close the outbound channel after all sends that channel are done.
func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	//start an output goroutine for each input channel in cs.
	//output copies values from c to out until c is closed. then calls wg.Done
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan int) {
			for n := range c {
				select {
				case out <- n:
				case <-done:
					return
				}
			}
			wg.Done()
		}(c)
	}

	//start a goroutine to close out once all the output goroutines are done.
	//this must start after the wg.Add call
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := sq(in)
	c2 := sq(in)

	// Consume the first value from output.
	done := make(chan struct{}, 2)
	defer close(done)
	out := merge(done, c1, c2)
	fmt.Println(<-out) // 4 or 9
}
