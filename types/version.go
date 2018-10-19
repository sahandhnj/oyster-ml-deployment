package types

import (
	"github.com/sahandhnj/apiclient/util"
)

type Version struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	VersionNumber    int    `json:"version_number"`
	ModelID          int    `json:"project_id"`
	Deployed         bool   `json:"deployed"`
	ImageTag         string `json:"image_tag"`
	ContainerId      string `json:"container_id"`
	RedisEnabled     bool   `json:"redis_enabled"`
	RedisContainerId string `json:"redis_container_id"`
	NetworkId        string `json:"network_id"`
	Port             int    `json:"port"`
}

const (
	RequirementsFilePath = "requirements.txt"
)

func NewVersion(versionNumber int, modelId int) (*Version, error) {
	uuid := util.UUID()

	v := Version{
		Name:          util.MinUUID(uuid),
		VersionNumber: versionNumber,
		ModelID:       modelId,
		Deployed:      false,
		RedisEnabled:  true,
	}

	return &v, nil
}
