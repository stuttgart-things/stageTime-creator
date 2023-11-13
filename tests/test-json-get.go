package main

import (
	"context"
	"fmt"
	"os"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/nitishm/go-rejson/v4"
)

var (
	redisServer        = os.Getenv("REDIS_SERVER")
	redisPort          = os.Getenv("REDIS_PORT")
	redisPassword      = os.Getenv("REDIS_PASSWORD")
	redisStream        = os.Getenv("REDIS_STREAM")
	ctx                = context.Background()
	revisionRunStageID = "385e2f8-0" //DEFAULT VALUE
)

func main() {

	// Check env vor given server port
	if os.Getenv("REVSIONRUN_STAGE_ID") != "" {
		revisionRunStageID = os.Getenv("REVSIONRUN_STAGE_ID")
	}

	// INITALIZE REDIS
	redisClient := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)
	redisJSONHandler := rejson.NewReJSONHandler()
	redisJSONHandler.SetGoRedisClient(redisClient)

	// GET ALL PIPELINERUN MANIFESTS FROM REDIS SET
	allPipelineRunNamesFromStage := sthingsCli.GetValuesFromRedisSet(redisClient, revisionRunStageID)
	fmt.Println("ALL PRS ON SET "+revisionRunStageID+" :", allPipelineRunNamesFromStage)

	// GET ALL PIPELINERUN MANIFESTS FROM REDIS SET
	for _, prName := range allPipelineRunNamesFromStage {
		manifestJSON := sthingsCli.GetRedisJSON(redisJSONHandler, prName)
		manifestYAML := sthingsCli.ConvertJSONToYAML(string(manifestJSON))
		fmt.Println("RENDERED YAML FOR PR " + prName + ":\n" + manifestYAML)
	}
}

func GetPipelineRunYAMLFromRedis(pipelineRunName string, redisJSONHandler *rejson.Handler) (pipelineRunYAML string) {
	pipelineRunJSON := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRunName)
	fmt.Println(string(pipelineRunJSON))
	pipelineRunYAML = sthingsCli.ConvertJSONToYAML(string(pipelineRunJSON))
	return pipelineRunYAML
}
