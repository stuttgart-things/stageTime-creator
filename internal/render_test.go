/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"testing"
)

const templateJobManifest = `
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .name }}
  namespace: {{ .namespace }}
`

const renderedJobManifest = `
apiVersion: batch/v1
kind: Job
metadata:
  name: test-job
  namespace: machine-shop
`

func TestRenderManifest(t *testing.T) {

	// MAPS HOLDING DATA
	redisValueData := make(map[string]string)
	templateValueData := make(map[string]interface{})

	// TEST DATA
	redisValueData["name"] = "test-job"
	redisValueData["namespace"] = "machine-shop"

	// ADD TEST DATA TO INTERFACE MAP
	for k, v := range redisValueData {
		templateValueData[k] = v
	}

	// TEST RENDER
	rendered := RenderManifest(templateValueData, templateJobManifest)

	// TEST OUTPUT
	fmt.Print(rendered)
	fmt.Print(renderedJobManifest)

	// COMPARE RESULT
	if rendered != renderedJobManifest {
		t.Errorf("expected: %s\ngot: %s", rendered, renderedJobManifest)
	}
}
