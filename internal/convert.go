/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"strings"
)

func validateCreateLoopValues(msgValues map[string]interface{}) map[string]interface{} {

	for key, value := range msgValues {

		if strings.Contains(key, "loop-") {
			delete(msgValues, key)
			_, key, _ = strings.Cut(key, "loop-")
			values := strings.Split(strings.TrimSpace(value.(string)), ";")

			msgValues[key] = map[string][]string{key: values}
		}

	}

	return msgValues
}
