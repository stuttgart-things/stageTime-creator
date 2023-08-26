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
	prList             = []string{"build-machineshop-image-1", "build-helm"}
)

func main() {

	pr := server.PipelineRun{
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

	// PUT PRS ON A LIST LOOP OVER AND USE pr.Name

	// CREATE REDIS CLIENT
	redisClient := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)

	// CREATE PIPELINERUNS ON REVISION RUN SET
	for _, pr := range prList {
		sthingsCli.AddValueToRedisSet(redisClient, revisionRunID, pr)
	}

	// PIPELINERUN ON REDIS JSON
	redisJSONHandler := rejson.NewReJSONHandler()
	redisJSONHandler.SetGoRedisClient(redisClient)
	sthingsCli.SetObjectToRedisJSON(redisJSONHandler, pr, pr.Name)
}
