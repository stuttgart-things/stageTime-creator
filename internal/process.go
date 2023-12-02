/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"

	"github.com/nitishm/go-rejson/v4"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stuttgart-things/redisqueue"
)

var (
	redisClient      = goredis.NewClient(&goredis.Options{Addr: redisServer + ":" + redisPort, Password: redisPassword, DB: 0})
	redisJSONHandler = rejson.NewReJSONHandler()
	targetNamespace  = "default"
)

func processStreams(msg *redisqueue.Message) error {

	log.Info("templatePath: ", templatePath)

	if msg.Values["stage"] != nil {
		fmt.Println("FOUND STAGE!")
		fmt.Println(msg.Values)
		redisJSONHandler.SetGoRedisClient(redisClient)

		revisionRunID := fmt.Sprintf("%v", msg.Values["revisionRunId"])
		fmt.Println(revisionRunID)
		allManifests := GetManifestFilesFromRedis(revisionRunID, redisJSONHandler)
		fmt.Println(allManifests)
		fmt.Println("PR0" + allManifests[0])
		targetNamespace, _ := sthingsBase.GetRegexSubMatch(allManifests[0], `namespace: "(.*?)"`)
		ApplyManifest(allManifests[0], targetNamespace)

	} else if msg.Values["template"] != nil {

		templateName := msg.Values["template"].(string)
		namespace = msg.Values["namespace"].(string)

		log.Info("templateName: ", templateName)
		log.Info("namespace: ", namespace)

		// verify values

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
