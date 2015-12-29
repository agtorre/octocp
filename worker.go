package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"
)

// NewWorker returns a generic worker type
func NewWorker(id int, workerQueue chan chan workRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan workRequest),
		WorkerQueue: workerQueue,
	}

	return worker
}

// Worker is the primary driver for executing various work types
type Worker struct {
	ID          int
	Work        chan workRequest
	WorkerQueue chan chan workRequest
}

// Start must be called by each worker to begin work
func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:

				switch work.Type {
				case workStat:
					w.Stat(work)
				case workReaddir:
					w.Readdir(work)
				}
				atomic.AddInt32(&busy, -1)

			case <-quitChan:
				log.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stat takes a an InPath, stats it, and if recursion is turned on sticks it
// on the work queue
func (w *Worker) Stat(wr workRequest) {
	info, err := os.Lstat(wr.InPath)
	if err != nil {
		log.Printf("error encountered: %s\n", err)
		stop()
		return
	}
	if Options.recursive && info.IsDir() {
		work := workRequest{
			Type:   workReaddir,
			InPath: wr.InPath,
		}
		workQueue <- work

	}
}

// Readdir takes an InPath, does a readdir, and puts abs paths to the files
// it find into the stat work queue
func (w *Worker) Readdir(wr workRequest) {
	files, err := ioutil.ReadDir(wr.InPath)
	if err != nil {
		log.Printf("error encountered: %s\n", err)
		stop()
		return
	}
	for _, f := range files {
		pRel := path.Join(wr.InPath, f.Name())
		pAbs, err := filepath.Abs(pRel)
		if err != nil {
			log.Printf("error encountered: %s\n", err)
			stop()
			return
		}
		work := workRequest{
			Type:   workStat,
			InPath: pAbs,
		}
		workQueue <- work
	}

}
