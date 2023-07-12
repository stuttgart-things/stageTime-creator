package main

// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	"github.com/redis/go-redis/v9"
// )

// func main() {
// 	hello := checkForRedisKV("hello")
// 	fmt.Println(hello)
// }

// func checkForRedisKV(name string) (jobIsFinished bool) {

// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_PORT"),
// 		Password: os.Getenv("REDIS_PASSWORD"),
// 		DB:       0,
// 	})

// 	// CHECK IF KEY EXISTS IN REDIS
// 	fmt.Println("CHECKING IF KEY " + name + " EXISTS..")
// 	keyExists, err := rdb.Exists(context.TODO(), name).Result()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// CHECK FOR VALUE/STATUS IN REDIS
// 	if keyExists == 1 {

// 		fmt.Println("KEY " + name + " EXISTS..CHECKING FOR IT'S VALUE")

// 		jobsStatus, err := rdb.Get(context.TODO(), name).Result()
// 		if err != nil {
// 			panic(err)
// 		}

// 		if jobsStatus == "finished" {
// 			jobIsFinished = true
// 		}

// 		fmt.Println("STATUS", jobsStatus)

// 	} else {

// 		fmt.Println("KEY " + name + " DOES NOT EXIST)")
// 	}

// 	return
// }
