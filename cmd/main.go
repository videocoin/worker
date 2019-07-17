package main

import (
	"log"
	"time"

	"github.com/videocoin/transcode"
)

func main() {
	err := transcode.Start()
	if err != nil {
		log.Println(err.Error())
	}
	time.Sleep(5 * time.Second)
}
