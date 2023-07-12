/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"testing"
)

func TestRenderManifest(t *testing.T) {

	for _, tc := range testsRender {

		// TEST RENDER
		rendered := RenderManifest(tc.testInput, tc.testTemplate)

		fmt.Println(rendered)
		fmt.Println(tc.want)

		if rendered != tc.want {
			t.Errorf("expected: %s\ngot: %s", rendered, renderedJobManifest)
		}

	}

}

const templateInventoryConfigMap = `
kind: ConfigMap
apiVersion: v1
metadata:
  name: ansible
  namespace: machine-shop
data:
  baseos-setup.yaml: |
    [{{ .groupName }}]
    {{ .hostName }}
`

const templateJobManifest = `
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .name }}
  namespace: {{ .namespace }}
`

const renderedInventoryConfigMap = `
kind: ConfigMap
apiVersion: v1
metadata:
  name: ansible
  namespace: machine-shop
data:
  baseos-setup.yaml: |
    [all]
    whatever.com
`

const renderedJobManifest = `
apiVersion: batch/v1
kind: Job
metadata:
  name: test-job
  namespace: machine-shop
`

type testRender struct {
	testTemplate string
	testInput    map[string]interface{}
	want         string
}

var (
	inventoryConfigMapValueData = map[string]interface{}{
		"groupName": "all",
		"hostName":  "whatever.com",
	}

	jobManifestValueData = map[string]interface{}{
		"name":      "test-job",
		"namespace": "machine-shop",
	}

	testsRender = []testRender{
		{testInput: jobManifestValueData, testTemplate: templateJobManifest, want: renderedJobManifest},
		{testInput: inventoryConfigMapValueData, testTemplate: templateInventoryConfigMap, want: renderedInventoryConfigMap},
	}
)
