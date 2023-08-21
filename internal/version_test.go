/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"testing"
)

const inventoryTemplate = `{{ range $name, $value := .inventory }}
[{{ $name }}]{{range $value }}
{{.}}{{end}}
{{ end }}`

var (
	ansibleInventory = map[string]interface{}{
		"all":                           "localhost",
		"loop-master":                   "rt.rancher.com;rt-2.rancher.com;rt-3.rancher.com",
		"loop-worker":                   "rt-4.rancher.com;rt-5.rancher.com",
		"merge-inventory;master;worker": "",
	}
)

func TestCreateLoopValues(t *testing.T) {
	loopableData, redisKey := validateCreateLoopValues(ansibleInventory)
	loopableData = validateMergeLoopValues(loopableData, redisKey)
	fmt.Println(loopableData)
	rendered := RenderManifest(loopableData, inventoryTemplate)
	fmt.Println(rendered)
}
