/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package main

import (
	"github.com/stuttgart-things/stageTime-creator/internal"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

var (
	log         = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
	logfilePath = "/tmp/stageTime-creator.log"
)

func main() {

	// PRINT BANNER + VERSION INFO
	internal.PrintBanner()
	log.Println("stageTime-creator started")

	// POLL FOR VALUES IN REDIS STREAMS
	internal.PollRedisStreams()
}
