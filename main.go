package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bosamatheus/gopher-download-manager/download"
)

func handleArgs() []string {
	if len(os.Args) < 4 {
		log.Fatal(fmt.Errorf("not enough arguments"))
	}
	return os.Args[1:]
}

func main() {
	args := handleArgs()
	d, err := download.New(args[0], args[1], args[2])
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	err = d.Do()
	if err != nil {
		log.Fatalf("an error has occurred while downloading the file: %s\n", err)
	}
	fmt.Printf("download completed in %v seconds\n", time.Since(start).Seconds())
}
