package pipe

import (
	"errors"
	"fmt"
)

/*
	Gets its input from a string channel, "processes" and sends itself to Processor out channel.
 */

type Ingest struct {
	Process
	In     <-chan string
	Out    chan Processor
	Result string
}

func (ig *Ingest) Run() (<-chan error, error) {
	if ig.In == nil {
		return nil, errors.New("nothing to ingest")
	}
	if ig.Out == nil {
		return nil, errors.New("requires a next process")
	}

	// We only need 1 error message in a buffer.
	errc := make(chan error, 1)

	go func() {
		defer close(errc)
		for {
			select {
			case in := <-ig.In:

				// If string is empty, skip it, but continue processing.
				if in == "" {
					fmt.Println("string is empty")
					break
				}

				ig.Result = in + " ingest "
				ig.Out <- ig
			case <-ig.context.Done():
				return
			}
		}

	}()

	return errc, nil
}
