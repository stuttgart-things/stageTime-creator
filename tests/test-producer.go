package main

import (
	"context"
	"fmt"
	"os"
	"time"

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
		{testValues: valuesConfigMap, testKey: "ConfigMap-ansible-inventory"},
		{testValues: ValuesJob, testKey: "Job-2023-07-02-configure-rke-node-19mv"},
	}
)

type test struct {
	testValues map[string]interface{}
	testKey    string
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

	// CHECK FOR VALUES IN REDIS
	for _, tc := range tests {

		fmt.Println(tc.testKey)

		retries := 0

		for range time.Tick(time.Second * 5) {

			if retries != 5 {

				retries = retries + 1
				if checkForRedisKV(tc.testKey, "created1") {
					break
				}

			} else {
				fmt.Println("not created!")

				break
			}

		}

	}

}

func checkForRedisKV(key, expectedValue string) (keyValueExists bool) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// CHECK IF KEY EXISTS IN REDIS
	fmt.Println("CHECKING IF KEY " + key + " EXISTS..")
	keyExists, err := rdb.Exists(context.TODO(), key).Result()
	if err != nil {
		panic(err)
	}

	// CHECK FOR VALUE/STATUS IN REDIS
	if keyExists == 1 {

		fmt.Println("KEY " + key + " EXISTS..CHECKING FOR IT'S VALUE")

		value, err := rdb.Get(context.TODO(), key).Result()
		if err != nil {
			panic(err)
		}

		if value == expectedValue {
			keyValueExists = true
		}

		fmt.Println("STATUS", value)

	} else {

		fmt.Println("KEY " + key + " DOES NOT EXIST)")
	}

	return
}
