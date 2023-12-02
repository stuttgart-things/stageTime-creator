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
	redisServer    = os.Getenv("REDIS_SERVER")
	redisPort      = os.Getenv("REDIS_PORT")
	redisPassword  = os.Getenv("REDIS_PASSWORD")
	redisStream    = os.Getenv("REDIS_STREAM")
	ctx            = context.Background()
	pipelineParams = make(map[string]string)
	resolverParams = make(map[string]string)

	listPipelineParams = make(map[string][]string)
	revisionRunStageID = "385e2f8-1"
	// pipelineWorkspaces   []server.Workspace
	volumeClaimTemplates []server.VolumeClaimTemplate

	volumeClaimTemplate = server.VolumeClaimTemplate{"shared-workspace", "openebs-hostpath", "ReadWriteOnce", "250Mi"}

	tektonPvc = server.Workspace{"volumeClaimTemplate", "volumeClaimTemplate", "codehub-ssh", "secretName"}
	workspace = server.Workspace{"source", "secret", "acr", "secretName"}
	// prList             = []string{"build-machineshop-image-1", "build-helm"}
	prs []string
)

func main() {

	fmt.Println("CONNECTED TO " + redisServer + ":" + redisPort)
	fmt.Println("STREAM " + redisStream)

	// GET REVISION RUN ID
	if os.Getenv("REVSIONRUN_STAGE_ID") != "" {
		revisionRunStageID = os.Getenv("REVSIONRUN_STAGE_ID")
	}

	// pipelineWorkspaces = append(pipelineWorkspaces, tektonPvc)
	volumeClaimTemplates = append(volumeClaimTemplates, volumeClaimTemplate)

	// var pipelineWorkspaces = append(pipelineWorkspaces, workspace)
	pr1 := server.PipelineRun{
		Name:                "pdns-cd43",
		RevisionRunAuthor:   "patrick.hermann",
		RevisionRunCreation: "23.1113.1007",
		RevisionRunCommitId: revisionRunStageID,
		RevisionRunRepoUrl:  "https://github.com/stuttgart-things/stuttgart-things.git",
		RevisionRunRepoName: "stuttgart-things",
		Namespace:           "stagetime-creator",
		// PipelineRef:         "create-kaniko-image",
		TimeoutPipeline: "0h12m0s",
		Params:          pipelineParams,
		ListParams:      listPipelineParams,
		ResolverParams:  resolverParams,
		Stage:           "0",
		NamePrefix:      "stagetime",
		NameSuffix:      "0",
		// Workspaces:           pipelineWorkspaces,
		VolumeClaimTemplates: volumeClaimTemplates,
	}

	resolverParams["url"] = "https://github.com/stuttgart-things/stuttgart-things.git"
	resolverParams["revision"] = "main"
	resolverParams["pathInRepo"] = "stageTime/pipelines/execute-ansible-playbooks.yaml"

	// SET KEY/VALUE PARAMETERS
	pipelineParams["ansibleWorkingImage"] = "eu.gcr.io/stuttgart-things/sthings-ansible:8.5.0"
	pipelineParams["createInventory"] = "false"
	pipelineParams["gitRepoUrl"] = "https://github.com/stuttgart-things/stuttgart-things.git"
	pipelineParams["gitRevision"] = "main"
	pipelineParams["gitWorkspaceSubdirectory"] = "/ansible/pdns"
	pipelineParams["vaultSecretName"] = "vault"
	pipelineParams["installExtraRoles"] = "true"

	// SET LIST PARAMETERS
	listPipelineParams["ansibleExtraRoles"] = []string{"https://github.com/stuttgart-things/install-configure-powerdns.git"}
	listPipelineParams["ansiblePlaybooks"] = []string{"ansible/playbooks/pdns-ingress-entry.yaml"}
	listPipelineParams["ansibleVarsFile"] = []string{"pdns_url+-https://pdns-pve.labul.sva.de:8443", "entry_zone+-sthings-pve.labul.sva.de.", "ip_address+-10.31.101.10", "hostname+-cd43"}
	listPipelineParams["ansibleVarsInventory"] = []string{"localhost"}

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
		fmt.Println(resourceName)
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
		"revisionRunId": revisionRunStageID,
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
