/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package main

import (
	"os"

	"github.com/stuttgart-things/stageTime-creator/internal"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

var (
	log         = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
	logfilePath = "/tmp/stageTime-creator.log"
	redisServer = os.Getenv("REDIS_SERVER")
	redisPort   = os.Getenv("REDIS_PORT")
	redisStream = os.Getenv("REDIS_STREAM")
)

func main() {

	// PRINT BANNER + VERSION INFO
	internal.PrintBanner()
	log.Info("STAGETIME-CREATOR STARTED")
	log.Info("REDIS SERVER " + redisServer)
	log.Info("REDIS PORT " + redisPort)
	log.Info("REDIS STREAM " + redisStream)

	// POLL FOR VALUES IN REDIS STREAMS
	internal.PollRedisStreams()
}
