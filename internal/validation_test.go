/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	pipelineRun = `
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  annotations:
    canfail: "true"
  labels:
    stagetime/author: patrick-hermann-sva
    stagetime/commit: 3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0
    stagetime/repo: stuttgart-things
    stagetime/stage: "1"
  name: st-1-simulate-stagetime-1713293c5a
  namespace: stagetime-tekton
spec:
  params:
  - name: gitRepoUrl
    value: https://github.com/stuttgart-things/stageTime-server.git
  - name: gitRevision
    value: main
  - name: gitWorkspaceSubdirectory
    value: stageTime
  - name: scriptPath
    value: tests/prime.sh
  - name: scriptTimeout
    value: 25s
  pipelineRef:
    params:
    - name: pathInRepo
      value: stageTime/pipelines/simulate-stagetime-pipelineruns.yaml
    - name: revision
      value: main
    - name: url
      value: https://github.com/stuttgart-things/stuttgart-things.git
    resolver: git
  taskRunTemplate:
    podTemplate:
      securityContext:
        fsGroup: 65532
  timeouts:
    pipeline: 0h30m0s
    tasks: 0h30m0s
  workspaces:
  - name: source
    volumeClaimTemplate:
      spec:
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: 20Mi
        storageClassName: openebs-hostpath
`
)

func TestValidateTemplateData(t *testing.T) {

	type test struct {
		scanText    string
		testPattern string
		want        []string
	}

	tests := []test{
		{scanText: "whatever {{ .Kind1 }}", testPattern: `\{\{(.*?)\}\}`, want: []string{"{{ .Kind1 }}"}},
		{scanText: "{{ .Name }}", testPattern: `\{\{(.*?)\}\}`, want: []string{"{{ .Name }}"}},
	}
	for _, tc := range tests {

		scanresult := GetAllRegexMatches(tc.scanText, tc.testPattern)
		fmt.Println(scanresult)
		fmt.Println(reflect.DeepEqual(scanresult, tc.want))
		if !reflect.DeepEqual(scanresult, tc.want) {
			t.Errorf("error")

		}
	}
}

func TestValidatePipelineRun(t *testing.T) {

	assert := assert.New(t)
	validPipelineRun, _ := ValidatePipelineRun(pipelineRun)

	assert.Equal(true, validPipelineRun)
}
