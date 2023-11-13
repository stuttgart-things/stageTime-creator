/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
	server "github.com/stuttgart-things/stageTime-server/server"
)

func GetPipelineRunFromRedis(redisJSONHandler *rejson.Handler, stageID string) (pipelineRun server.PipelineRun) {

	pipelineRunJSON, err := redis.Bytes(redisJSONHandler.JSONGet(stageID, "."))
	if err != nil {
		log.Fatalf("Failed to get pipelinRun from redis")
		return
	}

	pipelineRun = server.PipelineRun{}
	err = json.Unmarshal(pipelineRunJSON, &pipelineRun)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
		return
	}

	fmt.Printf("PipelineRun json read from redis : %#v\n", pipelineRun)

	return
}
