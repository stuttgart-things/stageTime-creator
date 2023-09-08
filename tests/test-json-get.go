package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/nitishm/go-rejson/v4"
	server "github.com/stuttgart-things/sweatShop-server/server"
)

var (
	redisServer        = os.Getenv("REDIS_SERVER")
	redisPort          = os.Getenv("REDIS_PORT")
	redisPassword      = os.Getenv("REDIS_PASSWORD")
	redisStream        = os.Getenv("REDIS_STREAM")
	ctx                = context.Background()
	revisionRunStageID = "7a6481c1-0"
)

func GetJSONFromRedis(pipelineRunName string, redisJSONHandler *rejson.Handler) (pipelineRun server.PipelineRun) {

	// pipelineRunJSON, err := redis.Bytes(redisJSONHandler.JSONGet(pipelineRunName, "."))
	// if err != nil {
	// 	log.Fatalf("Failed to JSONGet")
	// 	return
	// }

	pipelineRunJSON := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRunName)
	fmt.Println(string(pipelineRunJSON))

	// testJSON := `{"apiVersion":"tekton.dev/v1beta1", "kind":"PipelineRun", "metadata":{"name":"hello","namespace":"ansible"}}`
	// testJSON2 := `{"apiVersion":"tekton.dev/v1beta1","kind":"PipelineRun","metadata":{"labels":{"argocd.argoproj.io/instance":"tekton-runs","stagetime/author":"patrick-hermann-sva","stagetime/commit":"3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0","stagetime/repo":"stuttgart-things","tekton.dev/pipeline":"execute-ansible-smt40-rke2-15"},"name":"st-0-execute-ansible-smt40-rke2-15-0704523c5a","namespace":null},"spec":{"params":[{"name":"ansibleWorkingImage","value":"eu.gcr.io/stuttgart-things/sthings-ansible:8.3.0-8"},{"name":"createInventory","value":"true"},{"name":"gitRepoUrl","value":"git@github.com:stuttgart-things/stuttgart-things.git"},{"name":"gitRevision","value":"main"},{"name":"gitWorkspaceSubdirectory","value":"/ansible/base-os"},{"name":"installExtraRoles","value":"true"},{"name":"ansiblePlaybooks","value":["ansible/playbooks/prepare-env.yaml","ansible/playbooks/base-os.yaml"]}],"pipelineRef":{"name":"execute-ansible-smt40-rke2-15"},"serviceAccountName":"default","timeout":"1h","workspaces":[{"name":"ssh-credentials","secret":{"secretName":"codehub-ssh"}},{"name":"shared-workspace","persistentVolumeClaim":{"claimName":"ansible-tekton"}},{"name":"dockerconfig","secret":{"secretName":"scr-labda"}}]}}`

	// pipelineRunJSON1 := strings.ReplaceAll(string(pipelineRunJSON), "\\", "")
	// pipelineRunJSON1 = strings.TrimRight(pipelineRunJSON1, "\"")
	// pipelineRunJSON1 = strings.TrimLeft(pipelineRunJSON1, "\"")

	pipelineRunYAML := sthingsCli.ConvertJSONToYAML(string(pipelineRunJSON))
	// pipelineRunYAML2 := sthingsCli.ConvertJSONToYAML(pipelineRunJSON1)

	fmt.Println(pipelineRunYAML)
	// fmt.Println(pipelineRunYAML2)

	pipelineRun = server.PipelineRun{}
	err := json.Unmarshal(pipelineRunJSON, &pipelineRun)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal ", pipelineRunName)
		return
	}

	fmt.Printf("PipelineRun read from redis: %#v\n", pipelineRun)

	return
}

func main() {

	// INITALIZE REDIS
	redisClient := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)
	redisJSONHandler := rejson.NewReJSONHandler()
	redisJSONHandler.SetGoRedisClient(redisClient)

	GetJSONFromRedis("stageTime-server-test2", redisJSONHandler)

	// GET ALL PIPELINERUS FOR REVISION(ID)
	prs := sthingsCli.GetValuesFromRedisSet(redisClient, revisionRunStageID)

	for _, pr := range prs {

		fmt.Println(pr)
		GetJSONFromRedis(pr, redisJSONHandler)
		RenderManifest(pr, server.PipelineRunTemplate)

	}
}

func RenderManifest(resource interface{}, manifestTemplate string) string {

	var buf bytes.Buffer

	tmpl, err := template.New("manifest").Parse(manifestTemplate)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(&buf, resource)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
