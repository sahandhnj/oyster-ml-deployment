package types

import (
	"time"
)

type Req struct {
	ID           int       `json:"id" yaml:"id"`
	ModelId      int       `json:"model_id" yaml:"model_id"`
	VersionId    int       `json:"version_id" yaml:"version_id"`
	Result       bool      `json:"result" yaml:"result"`
	Time         time.Time `json:"time" yaml:"time"`
	ResponseTime int64     `json:"response_time" yaml:"response_time"`
}
