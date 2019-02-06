package main

import (
	"github.com/VideoCoin/transcode"
)

func main() {

	if err := transcode.Start(); err != nil {
		panic(err)
	}

}
