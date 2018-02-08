package main

import (
	"fmt"
)

// rightToLeft reads from the right channel, adds 1 and sends to the left.
func rightToLeft(left, right chan int) {
	// Blocks until the right most channel gets sent a message.
	left <- 1 + <-right
}

func main() {

	const n = 10

	// Reference to left most channel.
	leftmost := make(chan int)

	// We start of with one channel being the rightmost and the leftmost.
	right := leftmost
	left := leftmost

	// For as many channels as we want... (n)
	for i := 0; i < n; i++ {
		// Create a new channel to the "right"
		right = make(chan int)

		// Copy the value from the right to the left (and add 1)
		go rightToLeft(left, right)

		// The channel on the left of the to be created right (or the most right if the loop ends)
		left = right
	}

	// Kick it all off by writing a start int to the rightmost channel.
	// The will unblock all the blocked rightToLeft() calls.
	go func(c chan int) { c <- 0 }(right)

	// Block until left most channel has a result.
	fmt.Println(<-leftmost)

}
