package main

import (
	"fmt"
	"time"
	"math/rand"
)

// generator returns a receive only channel of strings.
func generator(name string) <-chan string {
	c := make(chan string)

	// Launch goroutine inside generator.
	go func(c chan string, n string) {
		for i := 0; ; i++ {
			// Write to channel in goroutine.
			c <- fmt.Sprintf("[%d] %s -- %s", i, n, time.Now())
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}(c, name)

	// Return channel to the caller.
	return c
}

func fanin(in1, in2 <-chan string) <-chan string {
	c := make(chan string)

	// Read from in1 and write to c
	go func() {
		for{
			select{
			case s := <-in1: c <- s
			case s := <-in2: c <- s
			}
		}
	}()

	return c
}

func main() {

	// Get the channel fed by the goroutine in the generator.
	chan1 := generator("Channel 1")
	chan2 := generator("Channel 2")
	fan := fanin(chan1, chan2)

	// Range over the channel. For unbuffered channel this synchronises.
	for m := range fan {
		fmt.Println(m)
	}

}
