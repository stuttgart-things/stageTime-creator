package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
	server "github.com/stuttgart-things/sweatShop-server/server"
)

var (
	redisServer   = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisStream   = os.Getenv("REDIS_STREAM")
	ctx           = context.Background()
)

func GetJSONFromRedis(redisJSONHandler *rejson.Handler) {

	studentJSON, err := redis.Bytes(redisJSONHandler.JSONGet("pipelineRun", "."))
	if err != nil {
		log.Fatalf("Failed to JSONGet")
		return
	}

	pipelineRun := server.PipelineRun{}
	err = json.Unmarshal(studentJSON, &pipelineRun)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}

	fmt.Printf("PipelineRun read from redis : %#v\n", pipelineRun)
}

func main() {

	// INITALIZE REDIS
	redisClient := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)

	prs := sthingsCli.GetValuesFromRedisSet(redisClient, "revisionrun-2432")
	fmt.Println(prs)

	redisJSONHandler := rejson.NewReJSONHandler()
	flag.Parse()

	redisJSONHandler.SetGoRedisClient(redisClient)

	GetJSONFromRedis(redisJSONHandler)

}
