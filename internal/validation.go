/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"regexp"
)

// func validateTemplateData() {

// 	scanText := "{{ .Kind }}"
// 	testPattern := `\{\{(.*?)\}\}`
// 	testOutput := GetAllRegexMatches(scanText, testPattern)
// 	fmt.Println(testOutput)

// }

func GetAllRegexMatches(scanText, regexPattern string) []string {

	re := regexp.MustCompile(regexPattern)
	fmt.Println(re)
	return re.FindAllString(scanText, -1)

}
