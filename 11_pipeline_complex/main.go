package main

/*
	Demonstrates a fairly complex pipeline.
	Note: The processes don't do much. Its just an example.
 */


import (
	"fmt"
	"github.com/rheinardkorf/go-concurrency/11_pipeline_complex/pipe"
	"context"
	"time"
	"io"
	"os"
)

var (

	// This is one source the messages can come from...
	polledSource     = generateMessages(context.Background())

	// But this is the one we need to watch and send to.
	// We could have multiple sources that all multiplex to this one.
	processQueue     = make(chan string)

	// This just returns a slice of processes to send to the pipe.
	processes        = initProcesses(processQueue, os.Stdout)

	// A kill signal for our main goroutine.
	terminateChannel = make(chan struct{})
)

// generateMessage uses a simple generator pattern to send messages to a string channel.
// The select statement may look odd here, but we're allowing our context to kill the
// nested goroutine.
func generateMessages(ctx context.Context) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		i := 0
		for {
			i += 1
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Println("Polling...")
			}

			out <- fmt.Sprintf("Msg %d:", i)
			time.Sleep(time.Second * 1)
		}
	}()

	return out
}

// initProcesses does as it says. It creates a number of processes for the pipeline.
// With each process we connect the "In" and "Out" channels to subsequent processes.
// For the Write process we're just being lazy. w is set to os.Stdout.
func initProcesses(source <-chan string, w io.Writer) []pipe.Processor {
	// Create and connect processes.
	ingest := &pipe.Ingest{
		In:  source,
		Out: make(chan pipe.Processor),
	}

	info := &pipe.Info{
		In:  ingest.Out,
		Out: make(chan pipe.Processor),
	}

	writer := &pipe.Write{
		In:     info.Out,
		Writer: w,
	}

	return []pipe.Processor{
		ingest,
		info,
		writer,
	}
}

func main() {

	// Create and run the pipe.
	go func() {
		// Create New Pipe with processes.
		pipeline := pipe.WithProcesses(processes...)

		err := pipeline.Run()

		// If our pipe doesn't run we need to send a signal on terminateChannel.
		// I personally like to use struct channels for this, but could be any other channel type.
		if err != nil {
			fmt.Println(err)
			terminateChannel <- struct{}{}
		}
	}()


	// We'll set a timeout for 5 minutes here, but adjust this accordingly to expectations.
	timeout := time.After(time.Minute * 5)

	// The run loop.
	for {
		select {
		case <-terminateChannel:
			return
		case msg := <-polledSource:
			// Send the message from the generator to the queue our pipe is really concerned with.
			processQueue <- msg
		case <-timeout:
			fmt.Println("5 minutes is a long time to wait for something to happen.")
			return
		}
	}
}
