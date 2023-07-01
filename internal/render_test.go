/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
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

	rendered := RenderManifest(Manifest{Name: "hello"}, renderedJobManifest)

	if rendered != renderedJobManifest {
		t.Errorf("expected: %s\ngot: %s", rendered, renderedJobManifest)
	}
}
