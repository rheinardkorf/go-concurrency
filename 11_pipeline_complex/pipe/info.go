package pipe

import (
	"errors"
)

/*
	Gets its input from a Processor channel, "processes" and sends itself to Processor out channel.
	Note: It also checks to see if the incoming message is an Ingest message.
 */

type Info struct {
	Process
	In     <-chan Processor
	Out    chan Processor
	Result string
}

func (i *Info) Run() (<-chan error, error) {
	if i.In == nil {
		return nil, errors.New("nothing to get info for")
	}
	if i.Out == nil {
		return nil, errors.New("requires a next process")
	}

	errc := make(chan error, 1)

	go func() {
		defer close(errc)
		for {
			select {
			case in := <-i.In:

				// Cast to *Ingest.
				obj, ok := in.(*Ingest)

				if ! ok {
					errc <- errors.New("expected an ingest process")
					break
				}

				i.Result = obj.Result + " info "
				i.Out <- i
			case <-i.context.Done():
				return
			}
		}

	}()

	return errc, nil
}
