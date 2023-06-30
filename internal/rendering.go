/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"bytes"
	"html/template"
)

type Manifest struct {
	Name string
}

const manifestTemplate = `
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .Name }}
  namespace: machine-shop
`

func RenderManifest(resource Manifest) string {

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
