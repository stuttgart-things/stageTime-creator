/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	sthingsK8s "github.com/stuttgart-things/sthingsK8s"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

func validateCreateLoopValues(msgValues map[string]interface{}) (map[string]interface{}, string) {

	// GENERATE KEY FOR THIS SET OPERATIONS
	redisKey, _ := time.Now().UTC().MarshalText()

	// CREATE REDIS CLIENT
	client := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)

	// CHECK VALUES FOR DATA BEGINING W/ LOOP
	for key, value := range msgValues {

		if strings.Contains(key, "loop-") {
			_, key, _ = strings.Cut(key, "loop-")
			values := strings.Split(strings.TrimSpace(value.(string)), ";")

			// REMOVE LOOP- ENTRY FROM VALUES MAP
			delete(msgValues, key)

			// ADD ALL LOOP DATA TO VALUES MAP
			msgValues[key] = map[string][]string{key: values}

			// ADD ALL LOOP DATA TO REDIS SETS (W/ THEIR KEYS)
			for _, value := range values {
				sthingsCli.AddValueToRedisSet(client, string(redisKey)+"-"+key, value)
			}

		}
	}

	return msgValues, string(redisKey)
}

func validateMergeLoopValues(msgValues map[string]interface{}, redisKey string) map[string]interface{} {

	// CREATE REDIS CLIENT
	client := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)

	// CHECK VALUES FOR DATA BEGINING W/ MERGE
	for key := range msgValues {

		var mergeKey string

		if strings.Contains(key, "merge-") {

			_, key, _ := strings.Cut(key, "merge-")

			mergedLoops := make(map[string][]string)

			// MERGE ALL LOOP DATA TO THE VALUES MAP BY THEIR MERGE KEY
			for i, key := range strings.Split(key, ";") {

				if i == 0 {
					mergeKey = key

				} else {
					mergedLoops[key] = sthingsCli.GetValuesFromRedisSet(client, redisKey+"-"+key)
				}

				msgValues[mergeKey] = mergedLoops
			}

		}
	}

	return msgValues
}

func GetAllRegexMatches(scanText, regexPattern string) []string {

	re := regexp.MustCompile(regexPattern)
	fmt.Println(re)
	return re.FindAllString(scanText, -1)

}

func ValidatePipelineRun(yamlPipelineRun string) (bool, error) {

	fmt.Println(yamlPipelineRun)

	// REWRITE PIPELINERUN FOR VALIDATION
	pr := strings.Split(yamlPipelineRun, "timeouts")
	paramsAsDefaults := strings.ReplaceAll("timeouts"+pr[1], "value:", "default:")
	rewrittenPipelineRun := pr[0] + paramsAsDefaults
	validPipelineRun, _, err := sthingsK8s.ConvertYAMLtoPipelineRun(rewrittenPipelineRun)

	return validPipelineRun, err
}
