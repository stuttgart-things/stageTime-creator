/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"

	"github.com/stuttgart-things/redisqueue"
)

func processStreams(msg *redisqueue.Message) error {

	log.Info("templatePath: ", templatePath)

	if msg.Values["stage"] != nil {
		fmt.Println("found stage!")
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
