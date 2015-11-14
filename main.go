package main

import (
	"flag"
	"fmt"
	"time"
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

type WorkKind string

const (
	workStat WorkKind = "stat"
)

type WorkRequest struct {
	Kind    WorkKind
	InPath  string
	OutPath string
}

var WorkQueue = make(chan WorkRequest, 100)

func processInput(args []string) {
	fmt.Println("processing input")

	for _, path := range args {
		work := WorkRequest{
			Kind:   workStat,
			InPath: path,
		}
		WorkQueue <- work
	}

}

func main() {
	flag.Parse()

	StartDispatcher(20)
	processInput(flag.Args())

	for {
		time.Sleep(time.Second * 1)
	}

}
