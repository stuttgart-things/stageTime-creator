/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package main

import (
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	"github.com/stuttgart-things/sweatShop-creator/internal"
)

var (
	log         = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
	logfilePath = "/tmp/sweatShop-creator.log"
)

func main() {

	// PRINT BANNER + VERSION INFO
	internal.PrintBanner()
	log.Println("sweatShop-creator started")

	// POLL FOR VALUES IN REDIS STREAMS
	internal.PollRedisStreams()
}
