/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"os"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsK8s "github.com/stuttgart-things/sthingsK8s"
)

func ApplyManifest(renderedManifest, namespace string) (manifestCreated bool) {

	clusterConfig, _ := sthingsK8s.GetKubeConfig(os.Getenv("KUBECONFIG"))

	kind, _ := sthingsBase.GetRegexSubMatch(renderedManifest, "kind:(.*)")
	resourceName, _ := sthingsBase.GetRegexSubMatch(renderedManifest, "name:(.*)")

	log.Info("TRYING TO APPLY " + kind + " W/ THE NAME " + resourceName)

	setRedisKV(kind+"-"+resourceName, "RENDERED")

	resourceCreated, resourceCreationError := sthingsK8s.CreateDynamicResourcesFromTemplate(clusterConfig, []byte(renderedManifest), namespace)
	if resourceCreationError != nil {
		log.Fatal(resourceCreationError)
	}

	if resourceCreated {
		setRedisKV(kind+"-"+resourceName, "CREATED")
		manifestCreated = true
	} else {
		manifestCreated = false
	}

	return
}
