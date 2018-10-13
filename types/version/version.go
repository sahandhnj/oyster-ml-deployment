package version

import (
	"github.com/sahandhnj/apiclient/filemanager"
	"github.com/sahandhnj/apiclient/types/model"
	"github.com/sahandhnj/apiclient/util"
)

type Version struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	DockerFile string `json:"dockerfile"`
	ModelID    string `json:"project_id"`
}

func NewVersion(model *model.Model) (*Version, error) {
	ID := util.UUID()

	v := Version{
		ID:         ID,
		Name:       util.MinUUID(ID),
		DockerFile: "Dockerfile-" + util.MinUUID(ID),
		ModelID:    model.Config.ID,
	}

	v.Apply(model)

	return &v, nil
}

func (v *Version) Apply(model *model.Model) error {
	fm, err := filemanager.NewFileStoreManager()
	if err != nil {
		return err
	}

	fm.CreateDirectoryInStore(v.Name)
	fm.CTarGz(v.Name+"/model.tar.gz", []string{model.Config.ModelPath}, false)

	return nil
}
