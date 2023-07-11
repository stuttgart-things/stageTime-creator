package main

import (
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/stuttgart-things/redisqueue"
)

var (
	redisServer   = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisStream   = os.Getenv("REDIS_STREAM")

	valuesConfigMap = map[string]interface{}{
		"template":  "inventory.gotmpl",
		"name":      "ansible-inventory",
		"namespace": "machine-shop",
		"groupName": "all",
		"hostName":  "whatever.com",
	}

	ValuesJob = map[string]interface{}{
		"template":  "ansible-job.yaml.gotmpl",
		"name":      "run-packer-rocky9",
		"namespace": "machine-shop",
	}

	tests = []test{
		{testValues: valuesConfigMap},
		{testValues: ValuesJob},
	}
)

type test struct {
	testValues map[string]interface{}
}

func main() {
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
