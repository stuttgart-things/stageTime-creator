/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"os"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsK8s "github.com/stuttgart-things/sthingsK8s"
)

func ApplyManifest(renderedManifest, namespace string) (manifestCreated bool) {

	clusterConfig, _ := sthingsK8s.GetKubeConfig(os.Getenv("KUBECONFIG"))

	kind, _ := sthingsBase.GetRegexSubMatch(renderedManifest, "kind:(.*)")
	resourceName, _ := sthingsBase.GetRegexSubMatch(renderedManifest, "name:(.*)")

	log.Info("trying to apply " + kind + " w/ the name " + resourceName)
	fmt.Println(renderedManifest)

	setRedisKV(kind+"-"+resourceName, "rendered")

	if sthingsK8s.CreateDynamicResourcesFromTemplate(clusterConfig, []byte(renderedManifest), namespace) {
		setRedisKV(kind+"-"+resourceName, "created")
		manifestCreated = true
	}

	return
}
