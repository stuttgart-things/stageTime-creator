/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"testing"

	sthingsK8s "github.com/stuttgart-things/sthingsK8s"

	"github.com/stretchr/testify/assert"
)

var (
	yamlPipelineRun = `
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: simulate-stagetime-pipelinerun-25
  namespace: tektoncd
spec:
  params:
    - name: gitRepoUrl
      default: 'https://github.com/stuttgart-things/stageTime-server.git'
    - name: gitRevision
      default: main
    - name: gitWorkspaceSubdirectory
      default: stageTime
    - name: scriptPath
      default: tests/prime.sh
    - name: scriptTimeout
      default: "15s"
  taskRunTemplate:
    podTemplate:
      securityContext:
        fsGroup: 65532
  workspaces:
    - name: source
      volumeClaimTemplate:
        spec:
          storageClassName: openebs-hostpath
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 20Mi
  pipelineRef:
    resolver: git
    params:
      - name: url
        value: "https://github.com/stuttgart-things/stuttgart-things.git"
      - name: revision
        value: main
      - name: pathInRepo
        value: stageTime/pipelines/simulate-stagetime-pipelineruns.yaml
`
)

func TestApplyManifest(t *testing.T) {

	assert := assert.New(t)

	validPipelineRun, _, _ := sthingsK8s.ConvertYAMLtoPipelineRun(yamlPipelineRun)
	assert.Equal(validPipelineRun, true)
}
