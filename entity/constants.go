package entity

import (
	"encoding/json"
)

type Result[ResultType interface{}] struct {
	Res ResultType
	Err error
}

func (res Result[ResultType]) ToJSON() string {
	resultInString, err0 := json.Marshal(res.Res)
	if err0 == nil {
		return string(resultInString)
	}
	return ""
}
