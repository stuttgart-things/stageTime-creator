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
	redisStream   = "q9:1"
)

func main() {

	internal.PrintBanner()
	template, templateFileExists := internal.ReadTemplateFromFilesystem(templatePath, templateName)

	fmt.Println(template)

	if templateFileExists {
		renderedManifest := internal.RenderManifest(internal.Manifest{Name: "dsfds"})
		fmt.Println(renderedManifest)
		internal.ApplyManifest(renderedManifest)
	} else {
		fmt.Println("TEMPLATE (PATH) NOT FOUND")
	}
}
