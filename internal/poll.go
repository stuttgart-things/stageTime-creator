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
	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

var (
	redisServer   = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisStream   = os.Getenv("REDIS_STREAM")
	templatePath  = os.Getenv("TEMPLATE_PATH")
	log           = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
	logfilePath   = "/tmp/sweatShop-creator.log"
	namespace     = "default" // just a default value
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

	log.Info("start polling stream", redisStream+" on "+redisServer+":"+redisPort)

	c.Run()

	log.Warn("polling stopped")

}

func processStreams(msg *redisqueue.Message) error {

	log.Info("templatePath: ", templatePath)
	fmt.Println(msg.Values["template"])

	if msg.Values["template"] != nil {

		templateName := msg.Values["template"].(string)
		namespace = msg.Values["namespace"].(string)

		log.Info("templateName: ", templateName)
		log.Info("namespace: ", namespace)

		// verify values

		template, templateFileExists := ReadTemplateFromFilesystem(templatePath, templateName)

		if templateFileExists {
			log.Info("template " + templateName + " imported")

			renderedManifest := RenderManifest(msg.Values, template)
			log.Info("rendered template: ", renderedManifest)

			ApplyManifest(renderedManifest, namespace)
		} else {
			log.Error("template " + templateName + " does not exist on filesystem")
		}

	} else {
		log.Error("templateName not defined in stream!")
	}

	return nil
}
