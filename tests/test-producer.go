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
)

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

	err2 := p.Enqueue(&redisqueue.Message{
		Stream: redisStream,
		Values: map[string]interface{}{
			"template": "job-template.yaml",
			"name":     "run-packer-rocky9",
		},
	})
	if err2 != nil {
		panic(err)
	}

}
