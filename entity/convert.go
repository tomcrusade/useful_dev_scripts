package entity

import (
	"encoding/json"
	"time"
)

func ConvertVultrTimeToDate(input string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, input)
	return parsedTime, err
}
func ConvertToJSON[inputType interface{}](input inputType) string {
	output, err0 := json.Marshal(input)
	if err0 == nil {
		return string(output)
	}
	return ""
}
