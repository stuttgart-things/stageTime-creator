package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
	goredis "github.com/redis/go-redis/v9"
	server "github.com/stuttgart-things/sweatShop-server/server"
)

var ctx = context.Background()

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
	var addr = flag.String("Server", "localhost:6379", "Redis server address")

	redisJSONHandler := rejson.NewReJSONHandler()
	flag.Parse()

	redisClient := goredis.NewClient(&goredis.Options{Addr: *addr, DB: 0})

	redisJSONHandler.SetGoRedisClient(redisClient)

	GetJSONFromRedis(redisJSONHandler)

}
