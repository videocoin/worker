package main

import (
	"log"
	"time"

	"github.com/VideoCoin/transcode"
)

func main() {

	for {
		if err := transcode.Start(); err != nil {
			log.Printf("failed to start, retrying...: %s", err.Error())
		}
		time.Sleep(1 * time.Minute)
	}

}
