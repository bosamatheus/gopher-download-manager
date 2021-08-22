package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bosamatheus/gopher-download-manager/download"
)

func main() {
	start := time.Now()
	d := download.New("https://golang.org/doc/gopher/run.png", "data/run.png", 10)
	err := d.Do()
	if err != nil {
		log.Fatalf("an error has occurred while downloading the file: %s\n", err)
	}
	fmt.Printf("download completed in %v seconds\n", time.Since(start).Seconds())
}
