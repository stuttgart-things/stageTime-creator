/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

// var (
// 	yamlPipelineRun = `
// apiVersion: tekton.dev/v1
// kind: PipelineRun
// metadata:
//   name: simulate-stagetime-pipelinerun-25
//   namespace: tektoncd
// spec:
//   pipelineRef:
//     resolver: git
//     params:
//       - name: url
//         value: "https://github.com/stuttgart-things/stuttgart-things.git"
//       - name: revision
//         value: main
//       - name: pathInRepo
//         value: stageTime/pipelines/simulate-stagetime-pipelineruns.yaml
//   params:
//     - name: gitRepoUrl
//       default: 'https://github.com/stuttgart-things/stageTime-server.git'
//     - name: gitRevision
//       default: main
//     - name: gitWorkspaceSubdirectory
//       default: stageTime
//     - name: scriptPath
//       default: tests/prime.sh
//     - name: scriptTimeout
//       default: "15s"
//   taskRunTemplate:
//     podTemplate:
//       securityContext:
//         fsGroup: 65532
//   workspaces:
//     - name: source
//       volumeClaimTemplate:
//         spec:
//           storageClassName: openebs-hostpath
//           accessModes:
//             - ReadWriteOnce
//           resources:
//             requests:
//               storage: 20Mi
// `
// )
