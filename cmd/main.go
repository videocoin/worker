package main

import (
	"log"
	"time"

	"github.com/VideoCoin/transcode"
)

func main() {

	for {
		if err := transcode.Start(); err != nil {
			log.Printf("no jobs found..")
		}
		time.Sleep(1 * time.Minute)
	}

}
