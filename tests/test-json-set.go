package main

import (
	"context"
	"os"

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
	pipelineWorkspaces []server.Workspace
	revisionRunID      = "7a6481c"
	// prList             = []string{"build-machineshop-image-1", "build-helm"}
	prs = []server.PipelineRun{}

	pr1 = server.PipelineRun{
		Name:                "build-machineshop-image-1",
		RevisionRunAuthor:   "patrick.hermann@sva.de",
		RevisionRunCreation: "pipelinerun.Name",
		RevisionRunCommitId: "pipelinerun.Name",
		RevisionRunRepoUrl:  "pipelinerun.Name",
		RevisionRunRepoName: "pipelinerun.Name",
		Namespace:           "pipelinerun.Name",
		PipelineRef:         "pipelinerun.Name",
		ServiceAccount:      "default",
		Timeout:             "1h",
		Params:              pipelineParams,
		Stage:               "pipelinerun.Name",
		Workspaces:          pipelineWorkspaces,
		NamePrefix:          "y",
		NameSuffix:          "pipelinerun.Name",
	}

	pr2 = server.PipelineRun{
		Name:                "package-machineshop-chart-1",
		RevisionRunAuthor:   "patrick.hermann@sva.de",
		RevisionRunCreation: "pipelinerun.Name",
		RevisionRunCommitId: "pipelinerun.Name",
		RevisionRunRepoUrl:  "pipelinerun.Name",
		RevisionRunRepoName: "pipelinerun.Name",
		Namespace:           "pipelinerun.Name",
		PipelineRef:         "pipelinerun.Name",
		ServiceAccount:      "default",
		Timeout:             "1h",
		Params:              pipelineParams,
		Stage:               "pipelinerun.Name",
		Workspaces:          pipelineWorkspaces,
		NamePrefix:          "y",
		NameSuffix:          "pipelinerun.Name",
	}
)

func main() {

	// PUT PRS ON A LIST
	prs = append(prs, pr1)
	prs = append(prs, pr2)

	// CREATE REDIS CLIENT
	redisClient := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)

	// PIPELINERUN ON REDIS JSON (LOOP OVER PRS AND USE pr.Name)
	redisJSONHandler := rejson.NewReJSONHandler()
	redisJSONHandler.SetGoRedisClient(redisClient)

	// CREATE PR REFERENCES (SET) AND OBJECTS (JSON) ON REDIS
	for _, pr := range prs {
		sthingsCli.AddValueToRedisSet(redisClient, revisionRunID, pr.Name)
		sthingsCli.SetObjectToRedisJSON(redisJSONHandler, pr, pr.Name)
	}
}
