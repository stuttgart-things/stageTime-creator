/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func validateCreateLoopValues(msgValues map[string]interface{}) map[string]interface{} {

	v, _ := time.Now().UTC().MarshalText()
	fmt.Println("REDIS KEY:", string(v))

	client := createRedisClient()

	for key, value := range msgValues {

		if strings.Contains(key, "loop-") {
			_, key, _ = strings.Cut(key, "loop-")
			values := strings.Split(strings.TrimSpace(value.(string)), ";")
			delete(msgValues, key)
			msgValues[key] = map[string][]string{key: values}

			for _, value := range values {
				AddValueToRedisSet(client, string(v)+"-"+key, value)
			}

		}
	}

	for key := range msgValues {

		var newName string

		if strings.Contains(key, "merge-") {

			_, key, _ := strings.Cut(key, "merge-")

			mergedLoops := make(map[string][]string)

			for i, key := range strings.Split(key, ";") {

				if i == 0 {
					newName = key

				} else {
					mergedLoops[key] = GetValuesFromRedisSet(client, string(v)+"-"+key)
				}

				msgValues[newName] = mergedLoops
			}

		}
	}

	return msgValues
}

func AddValueToRedisSet(redisClient *redis.Client, setKey, value string) (isSetValueunique bool) {

	if redisClient.SAdd(context.TODO(), setKey, value).Val() == 1 {
		isSetValueunique = true
	}

	return
}

func GetValuesFromRedisSet(redisClient *redis.Client, setKey string) (values []string) {

	values = redisClient.SMembers(context.TODO(), setKey).Val()

	return
}

func createRedisClient() (client *redis.Client) {

	client = redis.NewClient(&redis.Options{
		Addr:     redisServer + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})

	return
}
