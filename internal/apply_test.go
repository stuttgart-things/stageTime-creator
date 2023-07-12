/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"testing"
)

func TestApplyManifest(t *testing.T) {

	for _, tc := range testsApply {

		// TEST RENDER
		created := ApplyManifest(tc.renderedManifest, tc.namespace)

		fmt.Println(created)
		fmt.Println(tc.want)

		// if rendered != tc.want {
		// 	t.Errorf("expected: %s\ngot: %s", rendered, renderedJobManifest)
		// }

	}

}

type testApply struct {
	renderedManifest string
	namespace        string
	want             string
}

var (
	testsApply = []testApply{
		{renderedManifest: renderedInventoryConfigMap, namespace: "default", want: "created"},
	}
)
