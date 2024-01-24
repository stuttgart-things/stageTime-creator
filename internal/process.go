/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"encoding/json"
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
		revisionRunStatus := sthingsCli.GetRedisJSON(redisJSONHandler, revisionRunID+"-status")
		revisionRunFromRedis := server.RevisionRunStatus{}
		err := json.Unmarshal(revisionRunStatus, &revisionRunFromRedis)
		if err != nil {
			log.Fatalf("FAILED TO JSON UNMARSHAL REVISIONRUN STATUS")
		}
		server.PrintTable(revisionRunFromRedis)

		// PRINT STAGE STATUS
		stageStatus := sthingsCli.GetRedisJSON(redisJSONHandler, revisionRunID+stageNumber)
		stageStatusFromRedis := server.StageStatus{}
		err = json.Unmarshal(stageStatus, &stageStatusFromRedis)
		if err != nil {
			log.Fatalf("FAILED TO JSON UNMARSHAL STAGE STATUS")
		}
		server.PrintTable(stageStatusFromRedis)

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

	}

	// else if msg.Values["template"] != nil {

	// 	templateName := msg.Values["template"].(string)
	// 	namespace := msg.Values["namespace"].(string)

	// 	log.Info("templateName: ", templateName)
	// 	log.Info("namespace: ", namespace)

	// 	// VERIFY VALUES
	// 	template, templateFileExists := ReadTemplateFromFilesystem(templatePath, templateName)

	// 	if templateFileExists {

	// 		log.Info("template " + templateName + " imported")

	// 		log.Info("checking for loopable data..")
	// 		loopableData, redisKey := validateCreateLoopValues(msg.Values)
	// 		loopableData = validateMergeLoopValues(loopableData, redisKey)
	// 		fmt.Println(loopableData)

	// 		log.Info("rendering..")
	// 		renderedManifest := RenderManifest(msg.Values, template)
	// 		log.Info("rendered template: ", renderedManifest)

	// 		ApplyManifest(renderedManifest, namespace)

	// 	} else {
	// 		log.Error("template " + templateName + " does not exist on filesystem")
	// 	}

	// }

	return nil
}

func GetPipelineRunYAMLFromRedis(pipelineRunName string, redisJSONHandler *rejson.Handler) (pipelineRunYAML string) {

	pipelineRunJSON := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRunName)
	pipelineRunYAML = sthingsCli.ConvertJSONToYAML(string(pipelineRunJSON))

	return pipelineRunYAML

}
