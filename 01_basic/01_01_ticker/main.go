package main

import (
	"time"
	"fmt"
)

func ticker(iterations int, duration time.Duration) <-chan int {
	out := make(chan int)

	go func() {
		for i := 0; i < iterations; i++ {
			time.Sleep(duration)
			out <- i
		}
	}()

	return out
}

func main() {

	ticker1 := ticker(10, time.Second)
	ticker2 := ticker(5, time.Second*2)

	loop:
	for {
		select {
		case t1 := <-ticker1:
			fmt.Println("Ticker 1: ", t1)
		case t2 := <-ticker2:
			fmt.Println("Ticker 2: ", t2)
		case <-time.After(time.Second * 2): // Wait one second longer than expected next tick.
			fmt.Println("Done!")
			break loop
		}
	}
}
