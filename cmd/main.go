package main

import (
	"time"

	"github.com/VideoCoin/transcode"
)

func main() {

	for {
		transcode.Start()
		time.Sleep(1 * time.Minute)
	}

}
