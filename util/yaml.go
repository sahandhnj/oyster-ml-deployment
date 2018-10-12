package util

import (
	"gopkg.in/yaml.v2"
)

func MarshalYamlObject(object interface{}) ([]byte, error) {
	return yaml.Marshal(object)
}

func UnmarshalYamlObject(data []byte, object interface{}) error {
	return yaml.Unmarshal(data, object)
}
