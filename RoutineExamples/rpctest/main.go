package main

func main() {
	go startRPC()
	N := 1000
	mapChan := make(chan int, N)

	for i := 1; i < N; i ++{
		go func(i int) {
			
		}(i)
	}
}
