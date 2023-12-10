/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"encoding/json"
	"fmt"
	"os"

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

	log.Info("templatePath: ", templatePath)

	if msg.Values["stage"] != nil {

		log.Info(msg.Values)

		// TESTPRINT
		redisJSONHandler.SetGoRedisClient(redisClient)

		manifestJSON := sthingsCli.GetRedisJSON(redisJSONHandler, "3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0-1")

		revisionRunFromRedis := server.RevisionRunStatus{}
		err := json.Unmarshal(manifestJSON, &revisionRunFromRedis)
		if err != nil {
			log.Fatalf("Failed to JSON Unmarshal")
		}
		server.PrintTable(revisionRunFromRedis)
		// TESTPRINT

		revisionRunID := fmt.Sprintf("%v", msg.Values["revisionRunId"])
		fmt.Println(revisionRunID)
		allManifests := GetManifestFilesFromRedis(revisionRunID, redisJSONHandler)
		fmt.Println(allManifests)
		fmt.Println("PR0" + allManifests[0])
		ApplyManifest(allManifests[0], tektonNamespace)

	} else if msg.Values["template"] != nil {

		templateName := msg.Values["template"].(string)
		namespace := msg.Values["namespace"].(string)

		log.Info("templateName: ", templateName)
		log.Info("namespace: ", namespace)

		// VERIFY VALUES
		template, templateFileExists := ReadTemplateFromFilesystem(templatePath, templateName)

		if templateFileExists {

			log.Info("template " + templateName + " imported")

			log.Info("checking for loopable data..")
			loopableData, redisKey := validateCreateLoopValues(msg.Values)
			loopableData = validateMergeLoopValues(loopableData, redisKey)
			fmt.Println(loopableData)

			log.Info("rendering..")
			renderedManifest := RenderManifest(msg.Values, template)
			log.Info("rendered template: ", renderedManifest)

			ApplyManifest(renderedManifest, namespace)

		} else {
			log.Error("template " + templateName + " does not exist on filesystem")
		}

	} else {
		log.Error("templateName not defined in stream!")
	}

	return nil
}

func GetManifestFilesFromRedis(stageKey string, redisJSONHandler *rejson.Handler) (allManifests []string) {

	allManifestKeys := sthingsCli.GetValuesFromRedisSet(redisClient, stageKey)
	fmt.Println("ALL KEYS", allManifestKeys)

	for _, manifestKey := range allManifestKeys {

		// pr := GetPipelineRunYAMLFromRedis(manifestKey, redisJSONHandler)
		// fmt.Println(pr)

		manifestJSON := sthingsCli.GetRedisJSON(redisJSONHandler, manifestKey)
		// fmt.Println(sthingsCli.ConvertJSONToYAML(string(manifestJSON)))
		allManifests = append(allManifests, sthingsCli.ConvertJSONToYAML(string(manifestJSON)))
	}

	return

}

func GetPipelineRunYAMLFromRedis(pipelineRunName string, redisJSONHandler *rejson.Handler) (pipelineRunYAML string) {

	pipelineRunJSON := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRunName)
	pipelineRunYAML = sthingsCli.ConvertJSONToYAML(string(pipelineRunJSON))

	return pipelineRunYAML

}
