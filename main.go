/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package main

import (
	"github.com/stuttgart-things/sweatShop-creator/internal"
)

func main() {

	// PRINT BANNER + VERSION INFO
	internal.PrintBanner()

	// POLL FOR VALUES IN REDIS STREAMS
	internal.PollRedisStreams()
}
