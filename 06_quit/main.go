package main

import (
	"fmt"
	"time"
	"math/rand"
)

// generator returns a receive only channel of strings.
func generator(name string, quit chan<- bool) <-chan string {
	c := make(chan string)

	// Launch goroutine inside generator.
	go func(c chan string, n string) {
		for i := 0; ; i++ {
			// Write to channel in goroutine.
			c <- fmt.Sprintf("[%d] %s -- %s", i, n, time.Now())

			milliseconds := rand.Intn(1e3)

			// If its going to take more than 900 milliseconds, quit!
			if milliseconds >= 900 {
				quit <- true
			}

			time.Sleep(time.Duration(milliseconds) * time.Millisecond)
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

	// Channel for quitting.
	quit := make(chan bool)

	// Get the channel fed by the goroutine in the generator.
	chan1 := generator("Channel 1", quit)
	chan2 := generator("Channel 2", quit)
	fan := fanin(chan1, chan2)

	// Read from fan and/or quit.
	for{
		select {
		case s := <-fan:
			fmt.Println(s)
		case <-quit:
			fmt.Println("Somebody quit!")
			return
		}
	}

}
