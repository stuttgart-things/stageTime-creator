/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"testing"
)

const inventoryTemplate = `{{ range $name, $value := .master }}
[{{ $name }}]{{range $value }}
{{.}}{{end}}
{{ end }}`

var (
	ansibleInventory = map[string]interface{}{
		"all":         "localhost",
		"loop-master": "rt.rancher.com;rt-2.rancher.com;rt-3.rancher.com",
	}
)

func TestCreateLoopValues(t *testing.T) {
	loopableData := validateCreateLoopValues(ansibleInventory)
	fmt.Println(loopableData)
	rendered := RenderManifest(loopableData, inventoryTemplate)
	fmt.Println(rendered)
}
