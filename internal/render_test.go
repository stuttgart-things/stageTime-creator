/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"testing"
)

func TestRenderManifest(t *testing.T) {

	rendered := RenderManifest(playbookConfigMapValueData, templatePlaybookConfigMap)
	fmt.Println(rendered)

	for _, tc := range testsRender {

		// TEST RENDER
		rendered := RenderManifest(tc.testInput, tc.testTemplate)

		fmt.Println(rendered)
		fmt.Println(tc.want)

		if rendered != tc.want {
			t.Errorf("expected: %s\ngot: %s", rendered, tc.want)
		}

	}

}

const templateInventoryConfigMap = `
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ .name }}
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
  name: ansible-inventory
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

const templatePlaybookConfigMap = `
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ .name }}
  namespace: machine-shop
data:
  prepare-env.yaml: |
  ---
  - hosts: localhost
    become: false

    vars:
      home_dir: {{ "{{ lookup('env','HOME') }}" }}
      inv_path: {{ "{{ lookup('env','INV_PATH') }}" }}
      vault_approle_id: {{ "{{ lookup('env', 'VAULT_ROLE_ID') }}" }}
      vault_approle_secret: {{ "{{ lookup('env', 'VAULT_SECRET_ID') }}" }}
      vault_url: {{ "{{ lookup('env', 'VAULT_ADDR') }}" }}
`

type testRender struct {
	testTemplate string
	testInput    map[string]interface{}
	want         string
}

var (
	inventoryConfigMapValueData = map[string]interface{}{
		"name":      "ansible-inventory",
		"groupName": "all",
		"hostName":  "whatever.com",
	}

	jobManifestValueData = map[string]interface{}{
		"name":      "test-job",
		"namespace": "machine-shop",
	}

	playbookConfigMapValueData = map[string]interface{}{
		"name": "ansible-playbook",
	}

	testsRender = []testRender{
		{testInput: inventoryConfigMapValueData, testTemplate: templateInventoryConfigMap, want: renderedInventoryConfigMap},
		{testInput: jobManifestValueData, testTemplate: templateJobManifest, want: renderedJobManifest},
	}
)
