/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
	sthingsK8s "github.com/stuttgart-things/sthingsK8s"
)

var (
	pipelinRunJson = `{"apiVersion":"tekton.dev/v1","kind":"PipelineRun","metadata":{"name":"lint-machineshop-cli-47"},"spec":{"pipelineRef":{"resolver":"git","params":[{"name":"pathInRepo","value":"stageTime/pipelines/lint-golang-module.yaml"},{"name":"revision","value":"main"},{"name":"url","value":"https://github.com/stuttgart-things/stuttgart-things.git"}]},"workspaces":[{"name":"source","volumeClaimTemplate":{"spec":{"storageClassName":"nfs3-csi2","accessModes":["ReadWriteMany"],"resources":{"requests":{"storage":"1Gi"}}}}}],"params":[{"name":"gitRepoUrl","value":"https://github.com/stuttgart-things/machineShop.git"},{"name":"gitRevision","value":"main"},{"name":"gitWorkspaceSubdirectory","value":"machineShop"},{"name":"tokenSecretKey","value":"token"},{"name":"golintImage","value":"docker.io/golangci/golangci-lint:v1.54-alpine"},{"name":"ghImageURL","value":"maniator/gh:v2.35.0"}]}}`
)

func TestApplyManifest(t *testing.T) {

	assert := assert.New(t)

	namespace := "default"
	wantedKind := "PipelineRun"
	wantedResourceName := "lint-machineshop-cli-47"

	fmt.Println(pipelinRunJson)
	convertedYAMLpipelineRun := sthingsCli.ConvertJSONToYAML(pipelinRunJson)
	fmt.Println(convertedYAMLpipelineRun)

	kind, _ := sthingsBase.GetRegexSubMatch(convertedYAMLpipelineRun, "kind:(.*)")
	resourceName, _ := sthingsBase.GetRegexSubMatch(convertedYAMLpipelineRun, "name:(.*)")

	assert.Equal(kind, wantedKind)
	assert.Equal(resourceName, wantedResourceName)

	clusterConfig, _ := sthingsK8s.GetKubeConfig("/home/sthings/.kube/pve-cd43")
	resourceCreated, resourceCreationError := sthingsK8s.CreateDynamicResourcesFromTemplate(clusterConfig, []byte(convertedYAMLpipelineRun), namespace)
	if resourceCreationError != nil {
		log.Fatal(resourceCreationError)
	}

	assert.Equal(resourceCreated, true)

	// for _, tc := range testsApply {

	// 	// TEST RENDER
	// 	created := ApplyManifest(tc.renderedManifest, tc.namespace)

	// 	fmt.Println(created)
	// 	fmt.Println(tc.want)

	// 	// if rendered != tc.want {
	// 	// 	t.Errorf("expected: %s\ngot: %s", rendered, renderedJobManifest)
	// 	// }

	// }

}

type testApply struct {
	renderedManifest string
	namespace        string
	want             string
}

var (
	testsApply = []testApply{
		{renderedManifest: renderedInventoryConfigMap, namespace: "default", want: "created"},
	}
)
