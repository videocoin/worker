package main

import (
	"fmt"
	"time"

	"github.com/VideoCoin/transcode"
)

func main() {
	for {
		err := transcode.Start()
		if err != nil {
			fmt.Printf("error on start: %s retrying...", err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
	}

}
