/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package main

import (
	"fmt"
	"os"

	"github.com/stuttgart-things/sweatShop-creator/internal"
)

var (
	redisServer   = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisStream   = "q9:1"
)

func main() {

	internal.PrintBanner()

	hello := internal.RenderManifest(internal.Manifest{Name: "dsfds"})
	fmt.Println(hello)
}
