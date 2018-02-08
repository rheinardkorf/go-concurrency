package main

import (
	"fmt"
	"time"
)

// generator returns a receive only channel of strings.
func generator() <-chan string {
	c := make(chan string)

	// Launch goroutine inside generator.
	go func(c chan string) {
		for i := 0; ;i++ {
			// Write to channel in goroutine.
			c <- fmt.Sprintf("[%d] %s", i, time.Now())
			time.Sleep(time.Millisecond * 1000)
		}
	}(c)

	// Return channel to the caller.
	return c
}

func main() {

	// Get the channel fed by the goroutine in the generator.
	msg := generator()

	// Range over the channel. For unbuffered channel this synchronises.
	for m := range msg {
		fmt.Println(m)
	}

}
