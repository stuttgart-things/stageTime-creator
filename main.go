/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package main

import (
	"fmt"
	"os"

	"github.com/stuttgart-things/sweatShop-creator/internal"
)

var (
	templatePath  = os.Getenv("TEMPLATE_PATH")
	templateName  = os.Getenv("TEMPLATE_NAME")
	redisServer   = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisStream   = os.Getenv("REDIS_STREAM")
)

func main() {

	internal.PrintBanner()
	template, templateFileExists := internal.ReadTemplateFromFilesystem(templatePath, templateName)

	if templateFileExists {
		manifestValues := internal.Manifest{Name: "hello"}
		renderedManifest := internal.RenderManifest(manifestValues, template)
		fmt.Println(renderedManifest)
		internal.ApplyManifest(renderedManifest)
	} else {
		fmt.Println("TEMPLATE (PATH) NOT FOUND")
	}
}
