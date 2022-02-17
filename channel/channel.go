package main

import (
	"fmt"
	"time"
)

func Worker(i int) chan int {
	c := make(chan int)
	go func() {
		for {
			fmt.Printf("Worker %d received %c\n", i, <-c)
		}
	}()
	return c
}

func charDemo() {
	var channels [10]chan int
	for i := 0; i < 10; i++ {
		channels[i] = Worker(i)
	}
	for i := 0; i < 10; i++ {
		channels[i] <- 'a' + i
	}
	for i := 0; i < 10; i++ {
		channels[i] <- 'A' + i
	}
	time.Sleep(time.Millisecond)
}

func main() {
	charDemo()

}
