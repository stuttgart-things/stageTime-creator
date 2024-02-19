package main

import (
	"context"
	"fmt"
	"os"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsK8s "github.com/stuttgart-things/sthingsK8s"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/nitishm/go-rejson/v4"
)

var (
	redisServer            = os.Getenv("REDIS_SERVER")
	redisPort              = os.Getenv("REDIS_PORT")
	redisPassword          = os.Getenv("REDIS_PASSWORD")
	redisStream            = os.Getenv("REDIS_STREAM")
	ctx                    = context.Background()
	testRevisionRunStageID = "385e2f8-0" //DEFAULT VALUE
)

func main() {

	fmt.Println("CONNECTED TO " + redisServer + ":" + redisPort)
	fmt.Println("STREAM " + redisStream)

	// GET REVISION RUN ID
	if os.Getenv("REVSIONRUN_STAGE_ID") != "" {
		testRevisionRunStageID = os.Getenv("REVSIONRUN_STAGE_ID")
	}

	// INITALIZE REDIS
	redisClient := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)
	redisJSONHandler := rejson.NewReJSONHandler()
	redisJSONHandler.SetGoRedisClient(redisClient)

	// GET ALL PIPELINERUN MANIFESTS FROM REDIS SET
	allPipelineRunNamesFromStage := sthingsCli.GetValuesFromRedisSet(redisClient, testRevisionRunStageID)
	fmt.Println("ALL PRS ON SET "+testRevisionRunStageID+" :", allPipelineRunNamesFromStage)

	// GET ALL PIPELINERUN MANIFESTS FROM REDIS SET
	for _, prName := range allPipelineRunNamesFromStage {
		fmt.Println(prName)

		manifestJSON := sthingsCli.GetRedisJSON(redisJSONHandler, prName)
		manifestYAML := sthingsCli.ConvertJSONToYAML(string(manifestJSON))
		fmt.Println("RENDERED YAML FOR PR " + prName + ":\n" + manifestYAML)

		// CREATE PIPELINERUN ON THE CLUSTER
		clusterConfig, _ := sthingsK8s.GetKubeConfig(os.Getenv("KUBECONFIG"))

		kind, _ := sthingsBase.GetRegexSubMatch(manifestYAML, "kind:(.*)")
		resourceName, _ := sthingsBase.GetRegexSubMatch(manifestYAML, "name:(.*)")
		namespace, _ := sthingsBase.GetRegexSubMatch(manifestYAML, "namespace:(.*)")

		fmt.Println("APPLING " + kind + ": " + resourceName + "..")
		sthingsK8s.CreateDynamicResourcesFromTemplate(clusterConfig, []byte(manifestYAML), namespace)

	}
}
