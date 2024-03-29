/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/nitishm/go-rejson/v4"
	server "github.com/stuttgart-things/stageTime-server/server"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stuttgart-things/redisqueue"
)

var (
	redisClient      = goredis.NewClient(&goredis.Options{Addr: redisServer + ":" + redisPort, Password: redisPassword, DB: 0})
	redisJSONHandler = rejson.NewReJSONHandler()
	tektonNamespace  = os.Getenv("TEKTON_NAMESPACE")
)

func processStreams(msg *redisqueue.Message) error {
	redisJSONHandler.SetGoRedisClient(redisClient)

	if msg.Values["stage"] != nil {

		// PRINT STAGE DETAILS
		stage := fmt.Sprintf("%v", msg.Values["stage"])
		stageSplit := strings.Split(stage, "+")
		date := stageSplit[0]
		revisionRunID := stageSplit[1]
		stageNumber := stageSplit[2]

		log.Info("REVISIONRUN: ", revisionRunID)
		log.Info("STAGE DATE: ", date)
		log.Info("STAGE NUMBER: ", stageNumber)

		// PRINT REVISIONRUN STATUS
		revisionRunFromRedis := server.GetRevisionRunFromRedis(redisJSONHandler, revisionRunID+"-status", true)
		fmt.Println(revisionRunFromRedis)

		// PRINT STAGE STATUS
		stageStatusFromRedis := server.GetStageFromRedis(redisJSONHandler, revisionRunID+stageNumber, true)
		fmt.Println(stageStatusFromRedis)

		pipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageStatusFromRedis.StageID)
		log.Info("ALL PIPELINERUNS OF THIS STAGE: ", pipelineRuns)

		for _, pipelineRun := range pipelineRuns {
			log.Info("APPLING: ", pipelineRun)
			manifestJSON := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRun)
			pipelinRunYaml := sthingsCli.ConvertJSONToYAML(string(manifestJSON))

			// CHECK IF PIPELINERUN IS VALID
			validPipelineRun, prError := ValidatePipelineRun(pipelinRunYaml)

			// IF PIPELINERUN IS VALID - APPLY
			if validPipelineRun {
				ApplyManifest(pipelinRunYaml, tektonNamespace)
			} else {
				log.Error("PIPELINERUN INVALID: ", prError)
			}
		}

		// SET REVISIONRUN STATUS
		server.SetRevisionRunStatusInRedis(redisJSONHandler, revisionRunID+"-status", "REVISIONRUN UPDATED W/ CREATOR FOR STAGE: "+stageNumber, revisionRunFromRedis, true)

	}

	return nil
}

func GetPipelineRunYAMLFromRedis(pipelineRunName string, redisJSONHandler *rejson.Handler) (pipelineRunYAML string) {

	pipelineRunJSON := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRunName)
	pipelineRunYAML = sthingsCli.ConvertJSONToYAML(string(pipelineRunJSON))

	return pipelineRunYAML

}
