package main

import (
	"context"
	"fmt"
	"os"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/nitishm/go-rejson/v4"
)

var (
	redisServer   = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisStream   = os.Getenv("REDIS_STREAM")
	ctx           = context.Background()
	// revisionRunStageID = "7a6481c1-0"
	revisionRunStageID = "hello"
)

func GetPipelineRunYAMLFromRedis(pipelineRunName string, redisJSONHandler *rejson.Handler) (pipelineRunYAML string) {

	pipelineRunJSON := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRunName)
	fmt.Println(string(pipelineRunJSON))
	pipelineRunYAML = sthingsCli.ConvertJSONToYAML(string(pipelineRunJSON))

	return pipelineRunYAML

}

func main() {

	// INITALIZE REDIS
	redisClient := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)
	redisJSONHandler := rejson.NewReJSONHandler()
	redisJSONHandler.SetGoRedisClient(redisClient)
	// "st-0-build-kaniko-image-0940483c5a"
	pr := GetPipelineRunYAMLFromRedis("hello", redisJSONHandler)
	fmt.Println(pr)
	// GET ALL PIPELINERUS FOR REVISION(ID)
	// prs := sthingsCli.GetValuesFromRedisSet(redisClient, revisionRunStageID)

	// for _, pr := range prs {

	// 	fmt.Println(pr)
	// 	GetJSONFromRedis(pr, redisJSONHandler)
	// 	RenderManifest(pr, server.PipelineRunTemplate)

	// }
}
