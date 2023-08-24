package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	rejson "github.com/nitishm/go-rejson/v4"
	goredis "github.com/redis/go-redis/v9"
	server "github.com/stuttgart-things/sweatShop-server/server"
)

var ctx = context.Background()

func SetObjectToRedisJSON(redisJSONHandler *rejson.Handler, jsonObject interface{}, jsonKey string) {

	res, err := redisJSONHandler.JSONSet(jsonKey, ".", jsonObject)
	if err != nil {
		log.Fatalf("Failed to JSONSet")
		return
	}

	if res.(string) == "OK" {
		fmt.Printf("Success: %s\n", res)
	} else {
		fmt.Println("Failed to Set: ")
	}
}

func main() {

	//INITALIZE REDIS
	var addr = flag.String("Server", "127.0.0.1:28015", "Redis server address")

	pipelineParams := make(map[string]string)
	var pipelineWorkspaces []server.Workspace

	pr := server.PipelineRun{
		Name:                "pipelinerun.Name",
		RevisionRunAuthor:   "pipelinerun.Name",
		RevisionRunCreation: "pipelinerun.Name",
		RevisionRunCommitId: "pipelinerun.Name",
		RevisionRunRepoUrl:  "pipelinerun.Name",
		RevisionRunRepoName: "pipelinerun.Name",
		Namespace:           "pipelinerun.Name",
		PipelineRef:         "pipelinerun.Name",
		ServiceAccount:      "default",
		Timeout:             "1h",
		Params:              pipelineParams,
		Stage:               "pipelinerun.Name",
		Workspaces:          pipelineWorkspaces,
		NamePrefix:          "y",
		NameSuffix:          "pipelinerun.Name",
	}

	fmt.Println(pr)

	redisJSONHandler := rejson.NewReJSONHandler()
	flag.Parse()

	redisClient := goredis.NewClient(&goredis.Options{Addr: *addr, Password: os.Getenv("REDIS_PASSWORD"), DB: 0})

	redisJSONHandler.SetGoRedisClient(redisClient)

	SetObjectToRedisJSON(redisJSONHandler, pr, "pipelineRun")
}
