/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
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
	templatePath  = os.Getenv("TEMPLATE_PATH")
	templateName  = os.Getenv("TEMPLATE_NAME")
)

func PollRedisStreams() {

	c, err := redisqueue.NewConsumerWithOptions(&redisqueue.ConsumerOptions{
		VisibilityTimeout: 60 * time.Second,
		BlockingTimeout:   5 * time.Second,
		ReclaimInterval:   1 * time.Second,
		BufferSize:        100,
		Concurrency:       10,
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     redisServer + ":" + redisPort,
			Password: redisPassword,
			DB:       0,
		}),
	})

	if err != nil {
		panic(err)
	}

	c.Register(redisStream, processStreams)

	go func() {
		for err := range c.Errors {
			fmt.Printf("err: %+v\n", err)
		}
	}()

	fmt.Println("POLLING FOR REDIS STREAM " + redisStream + " ON " + redisServer + ":" + redisPort)

	c.Run()

	fmt.Println("POLLING FOR REDIS STREAM STOPPED")

}

func processStreams(msg *redisqueue.Message) error {

	// DEBUG OUTPUT MESSAGE
	fmt.Println("TEMPLATE-PATH", templatePath)
	fmt.Println("TEMPLATE", templateName)

	fmt.Printf("processing message: %v\n", msg.Values)
	// msg.Values["job"]

	// CHECK FOR TEMPLATE
	template, templateFileExists := ReadTemplateFromFilesystem(templatePath, templateName)

	if templateFileExists {
		manifestValues := Manifest{Name: "hello"}
		renderedManifest := RenderManifest(manifestValues, template)
		ApplyManifest(renderedManifest)

	} else {
		fmt.Println("TEMPLATE (PATH) NOT FOUND")
	}

	return nil
}
