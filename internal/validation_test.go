/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"reflect"
	"testing"
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
