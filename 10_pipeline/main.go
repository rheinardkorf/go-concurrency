package main

import (
	"context"
	"fmt"
	"errors"
	"sync"
	"time"
)

func PipeHead(ctx context.Context, stringc <-chan string) (<-chan string, <-chan error, error) {
	if stringc == nil {
		return nil, nil, errors.New("invalid input channel")
	}
	out := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)
		for {
			select {
			case in := <-stringc:
				if in == "" {
					errc <- errors.New("string is empty")
					break
				}
				out <- in + "...heads..."
			case <-ctx.Done():
				return
			}
		}

	}()

	return out, errc, nil
}

func PipeSink(ctx context.Context, stringc <-chan string) (<-chan error, error) {
	errc := make(chan error, 1)

	go func() {
		defer close(errc)
		for {
			select {
			case in := <-stringc:
				fmt.Println(in + "...sink..")
			case <-ctx.Done():
				return
			}
		}
	}()

	return errc, nil
}

func PipelineWait(errs []<-chan error) error {

	errc := PipelineMergeErrors(errs...)
	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

// MergeErrors merges multiple channels of errors.
// Based on https://blog.golang.org/pipelines.
// Based on https://medium.com/statuscode/pipeline-patterns-in-go-a37bb3a7e61d.
func PipelineMergeErrors(errs ...<-chan error) <-chan error {
	var wg sync.WaitGroup

	out := make(chan error, len(errs))

	// Create output function.
	output := func(ce <-chan error) {
		for err := range ce {
			out <- err
		}
		wg.Done()
	}

	wg.Add(len(errs))
	for _, err := range errs {
		go output(err)
	}

	// Drain the error channel.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func RunPipeline(in <-chan string) error {
	// Create new context with cancel function and defer.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a slice to contain all error channels.
	var allErrors []<-chan error

	// Pipeline head
	headc, errc, err := PipeHead(ctx, in)
	if err != nil {
		return err
	}
	allErrors = append(allErrors, errc)

	// Pipeline process

	// Pipeline sink
	errc, err = PipeSink(ctx, headc)
	if err != nil {
		return err
	}
	allErrors = append(allErrors, errc)

	fmt.Println("Running pipe.")

	return PipelineWait(allErrors)
}

func gen(quit <-chan bool, out chan string)  {

	go func(){
		defer close(out)
		i := 0
		for{
			i += 1
			select {
			case <-quit:
				return
			default:
				fmt.Println("Polling...")
			}

			out <- fmt.Sprintf("Msg %d:", i)
			time.Sleep(time.Second*1)
		}
	}()
}

func main() {

	done := make(chan bool)
	input := make(chan string)

	go func(){
		gen(done, input)
	}()

	go RunPipeline(input)

	time.Sleep(time.Second*10)
}
