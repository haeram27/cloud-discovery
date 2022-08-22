package util

import (
	"bytes"
	"cloudisc/apps"
	"encoding/json"
	"errors"

	"github.com/PaesslerAG/jsonpath"
)

func PrintPrettyJson(jsonBlob []byte) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, jsonBlob, "", "  ")
	if err != nil {
		apps.Logs.Error("JSON parse error: ", err.Error())
	} else {
		apps.Logs.Debug(prettyJSON.String())
	}
}

func PrettyJson(jsonBlob []byte) *bytes.Buffer {
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, jsonBlob, "", "  ")
	return &prettyJSON
}

func makeInterfaceArray(i interface{}) []interface{} {
	switch x := i.(type) {
	case []interface{}:
		return i.([]interface{})
	case interface{}:
		var arr []interface{}
		arr = append(arr, i)
		return arr
	default:
		_ = x
		apps.Logs.Error(errors.New("invalid type"))
		return []interface{}{}
	}
}

func JsonPath(jsonBlob []byte, path string) []interface{} {
	jsonStrt := interface{}(nil)
	err := json.Unmarshal(jsonBlob, &jsonStrt)
	if err != nil {
		apps.Logs.Error(err)
		return []interface{}{}
	}

	result, err := jsonpath.Get(path, jsonStrt)
	if err != nil {
		apps.Logs.Error(err)
		return []interface{}{}
	}

	/*
		jsonpath's result is
		interface{} (if size of result is 1)
		or
		[]interface{} (if size of result is 0 or size>1)

		make the result always as []interface{} for the convenient after process
	*/
	return makeInterfaceArray(result)
}
