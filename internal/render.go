/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"bytes"
	"html/template"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

func RenderManifest(resource map[string]interface{}, manifestTemplate string) string {

	var buf bytes.Buffer

	tmpl, err := template.New("manifest").Parse(manifestTemplate)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(&buf, resource)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func ReadTemplateFromFilesystem(templatePath, templateName string) (template string, templateFileExists bool) {

	templateFileExists, _ = sthingsBase.VerifyIfFileOrDirExists(templatePath+"/"+templateName, "file")

	if templateFileExists {

		template = sthingsBase.ReadFileToVariable(templatePath + "/" + templateName)

	}

	return
}
