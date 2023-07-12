/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func setRedisKV(key, value string) (kvSet bool) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisServer + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})

	err := rdb.Set(context.TODO(), key, value, 0).Err()

	if err != nil {
		panic(err)
	} else {
		kvSet = true
	}

	log.Info("KEY: " + key + " was set to VALUE: " + value + " on " + redisServer)

	return
}
