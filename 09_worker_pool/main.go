package main

import (
	"sync"
	"fmt"
	"time"
)

func worker(name string, tasks <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-tasks:
			if ! ok {
				return

			}

			// Delay to illustrate time passing...
			time.Sleep(time.Second)

			fmt.Printf("%s processed %s\n", name, task)
		}
	}
}

func pool(wg *sync.WaitGroup, workers int, tasks <-chan string) {

	name := "Worker"
	for i := 0; i < workers; i++ {
		go worker(fmt.Sprintf("%s %d", name, (i + 1)), tasks, wg)
	}

}

// generateTasks uses the generate pattern to fill a tasks channel.
// Note we need to drain the channel after the for loop or we'll
// end in a deadlock. We assume tasks are fixed for this example.
func generateTasks(count int) <-chan string {
	tasks := make(chan string)

	go func() {
		for i := 0; i < count; i++ {
			tasks <- fmt.Sprintf("Task %d", (i + 1))
		}
		close(tasks)
	}()

	return tasks
}

func main() {
	var wg sync.WaitGroup
	workerCount := 20
	taskCount := 200

	wg.Add(workerCount)

	tasks := generateTasks(taskCount)
	fmt.Println(tasks)
	go pool(&wg, workerCount, tasks)

	wg.Wait()
}
