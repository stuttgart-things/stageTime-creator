/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"testing"
)

const renderedJobManifest = `
apiVersion: batch/v1
kind: Job
metadata:
  name: hello
  namespace: machine-shop
`

func TestRenderManifest(t *testing.T) {

	hello := RenderManifest(Manifest{Name: "hello"})
	fmt.Println(hello)

	if hello != renderedJobManifest {
		t.Errorf("expected: %s\ngot: %s", hello, renderedJobManifest)
	}
}
