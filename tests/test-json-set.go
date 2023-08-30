package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	rejson "github.com/nitishm/go-rejson/v4"
	server "github.com/stuttgart-things/sweatShop-server/server"
)

var (
	redisServer        = os.Getenv("REDIS_SERVER")
	redisPort          = os.Getenv("REDIS_PORT")
	redisPassword      = os.Getenv("REDIS_PASSWORD")
	redisStream        = os.Getenv("REDIS_STREAM")
	ctx                = context.Background()
	pipelineParams     = make(map[string]string)
	listPipelineParams = make(map[string]string)
	revisionRunStageID = "7a6481c1-0"
	pipelineWorkspaces []server.Workspace
	tektonPvc          = server.Workspace{"ssh-credentials", "secret", "codehub-ssh", "secretName"}

	// prList             = []string{"build-machineshop-image-1", "build-helm"}
	prs = []server.PipelineRun{}

	pr2 = server.PipelineRun{
		Name:                "package-machineshop-chart-1",
		RevisionRunAuthor:   "patrick.hermann@sva.de",
		RevisionRunCreation: "pipelinerun.Name",
		RevisionRunCommitId: "pipelinerun.Name",
		RevisionRunRepoUrl:  "pipelinerun.Name",
		RevisionRunRepoName: "pipelinerun.Name",
		Namespace:           "tekton",
		PipelineRef:         "create-kaniko-image",
		ServiceAccount:      "default",
		Timeout:             "1h",
		Params:              pipelineParams,
		Stage:               "0",
		NameSuffix:          "0",
		NamePrefix:          "stageTime",
		ListParams:          listPipelineParams,
		Workspaces:          pipelineWorkspaces,
	}
)

func main() {

	pipelineWorkspaces = append(pipelineWorkspaces, tektonPvc)

	pr1 := server.PipelineRun{
		Name:                "build-machineshop-image-1",
		RevisionRunAuthor:   "patrick.hermann@sva.de",
		RevisionRunCreation: "pipelinerun.Name",
		RevisionRunCommitId: "pipelinerun.Name",
		RevisionRunRepoUrl:  "pipelinerun.Name",
		RevisionRunRepoName: "pipelinerun.Name",
		Namespace:           "tekton",
		PipelineRef:         "create-kaniko-image",
		ServiceAccount:      "default",
		Timeout:             "1h",
		Params:              pipelineParams,
		ListParams:          listPipelineParams,
		Stage:               "1",
		NamePrefix:          "stageTime",
		NameSuffix:          "1",
		Workspaces:          pipelineWorkspaces,
	}

	pipelineParams["image"] = "build-image"
	pipelineParams["tag"] = "123"
	listPipelineParams["gude"] = "123"
	// PUT PRS ON A LIST
	prs = append(prs, pr1)
	prs = append(prs, pr2)

	fmt.Println(pipelineWorkspaces)

	g := RenderManifest2(pr1, server.PipelineRunTemplate)
	fmt.Println(g)

	// CREATE REDIS CLIENT
	redisClient := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)

	// PIPELINERUN ON REDIS JSON (LOOP OVER PRS AND USE pr.Name)
	redisJSONHandler := rejson.NewReJSONHandler()
	redisJSONHandler.SetGoRedisClient(redisClient)

	// CREATE PR REFERENCES (SET) AND OBJECTS (JSON) ON REDIS
	for _, pr := range prs {
		sthingsCli.AddValueToRedisSet(redisClient, revisionRunStageID, pr.Name)
		sthingsCli.SetObjectToRedisJSON(redisJSONHandler, pr, pr.Name)
		fmt.Println(pr)
	}
}

func RenderManifest2(resource interface{}, manifestTemplate string) string {

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
