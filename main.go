package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"sync"
)

var Options struct {
	verbose   bool
	recursive bool
}

func config() {
	const (
		verboseDefault   = false
		verboseUsage     = "Display additional output"
		recursiveDefault = false
		recursiveUsage   = "Recursively traverse all directories"
	)
	flag.BoolVar(&Options.verbose, "verbose", verboseDefault, verboseUsage)
	flag.BoolVar(&Options.verbose, "v", verboseDefault, verboseUsage+" (shorthand)")
	flag.BoolVar(&Options.recursive, "recursive", recursiveDefault, recursiveUsage)
	flag.BoolVar(&Options.recursive, "r", recursiveDefault, recursiveUsage+" (shorthand)")

}

func init() {
	config()
}

//Begin magic
type WorkType string

const (
	workStat    WorkType = "stat"
	workReaddir          = "readdir"
)

type workRequest struct {
	Type    WorkType
	InPath  string
	OutPath string
}

var workQueue = make(chan workRequest, 100)

func processInput(args []string) {
	fmt.Println("processing input")

	for _, path := range args {
		pAbs, err := filepath.Abs(path)
		if err != nil {
			log.Panicf("error encountered: %\n", err)
		}
		work := workRequest{
			Type:   workStat,
			InPath: pAbs,
		}
		workQueue <- work
	}

}

func main() {
	flag.Parse()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	StartDispatcher(5, wg)
	processInput(flag.Args())
	wg.Wait()
}
