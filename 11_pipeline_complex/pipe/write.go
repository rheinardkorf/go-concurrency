package pipe

import (
	"errors"
	"io"
)

/*
	Gets its input from a Processor channel, "processes" and sends itself to Processor out channel.
	Note: It also checks to see if the incoming message is an Info message.
 */

type Write struct {
	Process
	In     <-chan Processor
	Writer io.Writer
	Results int
}

func (w *Write) Run() (<-chan error, error) {
	if w.In == nil {
		return nil, errors.New("nothing to get info for")
	}

	errc := make(chan error, 1)

	go func() {
		defer close(errc)
		for {
			select {
			case in := <-w.In:

				// Cast to *Ingest.
				obj, ok := in.(*Info)

				if ! ok {
					errc <- errors.New("expected an info process")
					break
				}

				if w.Writer == nil {
					errc <- errors.New("no 'Writer' defined")
					break
				}

				message := obj.Result + " writer \n"

				var err error

				w.Results, err = w.Writer.Write([]byte(message))

				if err != nil {
					errc <- err
					break
				}
			case <-w.context.Done():
				return
			}
		}

	}()

	return errc, nil
}
