package main

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/stuttgart-things/redisqueue"
)

var (
	redisServer   = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisStream   = os.Getenv("REDIS_STREAM")

	tests = []test{
		// {testValues: valuesConfigMap, testKey: "ConfigMap-ansible-inventory"},
		// {testValues: ValuesJob, testKey: "Job-2023-07-02-configure-rke-node-19mv"},
		// {testValues: ValuesPlaybook, testKey: "tbd"},
		{testValues: ValuesStage, testKey: "tbd"},
	}

	// valuesConfigMap = map[string]interface{}{
	// 	"template":                      "inventory.gotmpl",
	// 	"name":                          "ansible-inventory",
	// 	"namespace":                     "machine-shop",
	// 	"all":                           "localhost",
	// 	"loop-master":                   "rt.rancher.com;rt-2.rancher.com;rt-3.rancher.com",
	// 	"loop-worker":                   "rt-4.rancher.com;rt-5.rancher.com",
	// 	"merge-inventory;master;worker": "",
	// }

	ValuesPlaybook = map[string]interface{}{
		"template":  "playbook.gotmpl",
		"name":      "play",
		"namespace": "machine-shop",
	}

	ValuesStage = map[string]interface{}{
		"stage":         "stage0",
		"kind":          "pipelinRun",
		"revisionRunId": "7a6481c1-0",
	}

	// ValuesJob = map[string]interface{}{
	// 	"template":  "playbook.gotmpl",
	// 	"name":      "playy",
	// 	"namespace": "machine-shop",
	// }

)

type test struct {
	testValues map[string]interface{}
	testKey    string
}

func main() {

	fmt.Println("REDIS-SERVER: " + redisServer + ":" + redisPort)
	fmt.Println("REDIS-STREAM: " + redisStream)

	// CREATE RESOURCES IN REDIS
	p, err := redisqueue.NewProducerWithOptions(&redisqueue.ProducerOptions{
		ApproximateMaxLength: true,
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     redisServer + ":" + redisPort,
			Password: redisPassword,
			DB:       0,
		}),
	})

	if err != nil {
		panic(err)
	}

	// CREATE RESOURCES IN REDIS
	for _, tc := range tests {

		err2 := p.Enqueue(&redisqueue.Message{
			Stream: redisStream,
			Values: tc.testValues,
		})

		if err2 != nil {
			panic(err)
		}

	}

}
