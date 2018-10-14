package types

import (
	"github.com/sahandhnj/apiclient/util"
)

type Version struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	VersionNumber int    `json:"version_number"`
	ModelID       int    `json:"project_id"`
	Deployed      bool   `json:"deployed"`
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
	}

	return &v, nil
}
