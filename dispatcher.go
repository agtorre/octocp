package main

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var workerQueue chan chan workRequest
var quitChan chan bool
var busy int32 

func stop() {
	close(quitChan)
}

// StartDispatcher launches a dispatcher who's role is to dispatch work, but also
// keep track of the state of the world. When there's no more work and workers are done
// processing work, the Dispatcher terminates the program
func StartDispatcher(nworkers int, wg *sync.WaitGroup) {
	// First, initialize the channel we are going to but the workers' work channels into.
	workerQueue = make(chan chan workRequest, nworkers)
	quitChan = make(chan bool)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		log.Println("Starting worker", i)
		worker := NewWorker(i, workerQueue)
		worker.Start()
	}

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {

			select {
			case work := <-workQueue:
				log.Printf("Received work request: %s\n", work)
				atomic.AddInt32(&busy, 1)
				go func() {
					worker := <-workerQueue

					log.Println("Dispatching work request")
					worker <- work
				}()
			case <-time.After(time.Second * 1):
				if atomic.LoadInt32(&busy) == 0 && len(workQueue) == 0 {
					close(quitChan)
					return
				}
			}
		}
	}(wg)
}
