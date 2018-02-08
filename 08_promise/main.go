package main

import (
	"time"
	"fmt"
)

// promiseChannel demonstrates value that we will have in the future.
// This is ideal for Async requests, but here we're just using a timer.
func promiseChannel(input string) <-chan string {

	c := make(chan string)

	// Use a goroutine to get our future data and return it via the channel.
	go func(){

		// E.g. In the real world, do an http requests here.
		time.Sleep(time.Second * 5)

		c <- "Hi " + input + ". This is your promised message."
	}()

	return c
}

func main() {

	// A new promised future value.
	promise := promiseChannel("Promise 1")

	for i:=0; i < 10; i++ {
		fmt.Println("Item: ", i)
	}

	fmt.Println("Lots of things can happen until we need our promised data, but now we have to wait...")

	// At this point, we have to have our promised data.
	promisedMessage := <-promise

	fmt.Println(promisedMessage)
}