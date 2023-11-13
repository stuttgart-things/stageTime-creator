package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	rejson "github.com/nitishm/go-rejson/v4"
	server "github.com/stuttgart-things/stageTime-server/server"
)

var (
	redisServer        = os.Getenv("REDIS_SERVER")
	redisPort          = os.Getenv("REDIS_PORT")
	redisPassword      = os.Getenv("REDIS_PASSWORD")
	redisStream        = os.Getenv("REDIS_STREAM")
	ctx                = context.Background()
	pipelineParams     = make(map[string]string)
	listPipelineParams = make(map[string][]string)
	revisionRunStageID = "7a6481c1-0"
	pipelineWorkspaces []server.Workspace

	tektonPvc = server.Workspace{"ssh-credentials", "secret", "codehub-ssh", "secretName"}
	workspace = server.Workspace{"source", "secret", "acr", "secretName"}
	// prList             = []string{"build-machineshop-image-1", "build-helm"}
	prs []string
)

func main() {

	pipelineWorkspaces = append(pipelineWorkspaces, tektonPvc)

	var pipelineWorkspaces = append(pipelineWorkspaces, workspace)
	pr1 := server.PipelineRun{
		Name:                "build-machineshop-image-0",
		RevisionRunAuthor:   "patrick.hermann@sva.de",
		RevisionRunCreation: "23.1113.1007",
		RevisionRunCommitId: "385e2f8",
		RevisionRunRepoUrl:  "https://github.com/stuttgart-things/stuttgart-things.git",
		RevisionRunRepoName: "stuttgart-things",
		Namespace:           "tekton",
		PipelineRef:         "create-kaniko-image",
		ServiceAccount:      "default",
		Timeout:             "1h",
		Params:              pipelineParams,
		ListParams:          listPipelineParams,
		Stage:               "0",
		NamePrefix:          "stageTime",
		NameSuffix:          "0",
		Workspaces:          pipelineWorkspaces,
	}

	// SET PARAMETERS
	pipelineParams["image"] = "build-image"
	pipelineParams["tag"] = "123"
	listPipelineParams["gude"] = []string{"123"}

	// TEST RENDER
	renderedPr := RenderPipelineRun(pr1, server.PipelineRunTemplate)

	// PUT PRS ON A LIST
	prs = append(prs, renderedPr)

	// CREATE REDIS CLIENT
	redisClient := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)

	// PIPELINERUN ON REDIS JSON (LOOP OVER PRS AND USE pr.Name)
	redisJSONHandler := rejson.NewReJSONHandler()
	redisJSONHandler.SetGoRedisClient(redisClient)

	// CREATE PR REFERENCES (SET) AND OBJECTS (JSON) ON REDIS
	for _, pr := range prs {
		resourceName, _ := sthingsBase.GetRegexSubMatch(pr, `name: "(.*?)"`)
		sthingsCli.AddValueToRedisSet(redisClient, revisionRunStageID, resourceName)
		sthingsCli.SetRedisJSON(redisJSONHandler, sthingsCli.ConvertYAMLToJSON(pr), resourceName)
		fmt.Println(pr)
		fmt.Println("STORED PR " + resourceName + " ON SET " + revisionRunStageID)
		fmt.Println("STORED PR " + resourceName + " AS JSON ON " + resourceName)
	}

	// CREATE DATA ON REDIS STREAMS
	ValuesStage := map[string]interface{}{
		"stage":         "stage0",
		"kind":          "pipelinRun",
		"revisionRunId": "7a6481c1-0",
	}

	sthingsCli.EnqueueDataInRedisStreams(redisServer+":"+redisPort, redisPassword, redisStream, ValuesStage)
	fmt.Println("STORED PR DATA ON STREAM", redisStream)
}

func RenderPipelineRun(resource interface{}, manifestTemplate string) string {

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
