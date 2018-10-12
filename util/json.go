package util

import (
	"encoding/json"
)

func MarshalJsonObject(object interface{}) ([]byte, error) {
	return json.Marshal(object)
}

func UnmarshalJsonObject(data []byte, object interface{}) error {
	return json.Unmarshal(data, object)
}
