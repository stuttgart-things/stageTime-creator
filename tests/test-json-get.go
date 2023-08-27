package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
	server "github.com/stuttgart-things/sweatShop-server/server"
)

var (
	redisServer        = os.Getenv("REDIS_SERVER")
	redisPort          = os.Getenv("REDIS_PORT")
	redisPassword      = os.Getenv("REDIS_PASSWORD")
	redisStream        = os.Getenv("REDIS_STREAM")
	ctx                = context.Background()
	revisionRunStageID = "7a6481c1-0"
)

func GetJSONFromRedis(pipelineRunName string, redisJSONHandler *rejson.Handler) (pipelineRun server.PipelineRun) {

	pipelineRunJSON, err := redis.Bytes(redisJSONHandler.JSONGet(pipelineRunName, "."))
	if err != nil {
		log.Fatalf("Failed to JSONGet")
		return
	}

	pipelineRun = server.PipelineRun{}
	err = json.Unmarshal(pipelineRunJSON, &pipelineRun)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal ", pipelineRunName)
		return
	}

	fmt.Printf("PipelineRun read from redis: %#v\n", pipelineRun)

	return
}

func main() {

	// INITALIZE REDIS
	redisClient := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)
	redisJSONHandler := rejson.NewReJSONHandler()
	redisJSONHandler.SetGoRedisClient(redisClient)

	// GET ALL PIPELINERUS FOR REVISION(ID)
	prs := sthingsCli.GetValuesFromRedisSet(redisClient, revisionRunStageID)

	for _, pr := range prs {

		fmt.Println(pr)
		GetJSONFromRedis(pr, redisJSONHandler)

	}
}
