/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"os"

	sthingsK8s "github.com/stuttgart-things/sthingsK8s"
)

func ApplyManifest(renderedManifest string) {

	clusterConfig, _ := sthingsK8s.GetKubeConfig(os.Getenv("KUBECONFIG"))

	// DEBUG
	ns := sthingsK8s.GetK8sNamespaces(clusterConfig)
	fmt.Println("FOUND NAMESAPCES", ns)

	sthingsK8s.CreateDynamicResourcesFromTemplate(clusterConfig, []byte(renderedManifest), "default")

}
