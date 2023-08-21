/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"flag"
	"fmt"

	"github.com/nitishm/go-rejson/v4"
	server "github.com/stuttgart-things/sweatShop-server/server"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stuttgart-things/redisqueue"
)

func processStreams(msg *redisqueue.Message) error {

	log.Info("templatePath: ", templatePath)

	if msg.Values["stage"] != nil {
		fmt.Println("found stage!")
		var addr = flag.String("Server", "10.100.136.56:31868", "Redis server address")
		redisJSONHandler := rejson.NewReJSONHandler()
		flag.Parse()
		redisClient := goredis.NewClient(&goredis.Options{Addr: *addr, DB: 0})
		redisJSONHandler.SetGoRedisClient(redisClient)
		pr := GetPipelineRunFromRedis(redisJSONHandler, "pipelineRun")

		renderedManifest := RenderManifest(pr, server.PipelineRunTemplate)
		log.Info("rendered template: ", renderedManifest)

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
