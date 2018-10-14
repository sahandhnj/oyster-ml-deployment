package version

import (
	"path"

	"github.com/sahandhnj/apiclient/filemanager"
	"github.com/sahandhnj/apiclient/types/model"
	"github.com/sahandhnj/apiclient/util"
)

type Version struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	DockerFile string `json:"dockerfile"`
	ModelID    int    `json:"project_id"`
}

const (
	RequirementsFilePath = "requirements.txt"
)

func NewVersion(ID int, model *model.Model) (*Version, error) {
	uuid := util.UUID()

	v := Version{
		ID:         ID,
		Name:       util.MinUUID(uuid),
		DockerFile: "Dockerfile-" + util.MinUUID(uuid),
		ModelID:    model.ID,
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
	fm.CTarGz(path.Join(v.Name, "model.tar.gz"), []string{model.ModelPath}, false)
	fm.CopyToStore(path.Join(model.ModelPath, "requirements.txt"), path.Join(v.Name, "requirements.txt"))
	v.createDockerFile(fm)

	return nil
}

func (v *Version) createDockerFile(fm *filemanager.FileStoreManager) {
	docker_file_static = docker_file_static + "RUN pip install --user " + fm.ReadRQLineByLine(path.Join(v.Name, "requirements.txt"))
	fm.WriteToFile(path.Join(v.Name, "Dockerfile"), docker_file_static)
	fm.WriteToFile(path.Join(v.Name, ".dockerignore"), "")
}

var docker_file_static = `FROM ubuntu:18.04
ENV MODELPATH /src

RUN apt-get update && apt-get install -y software-properties-common

RUN add-apt-repository ppa:jonathonf/python-3.6 
RUN apt-get update && apt-get install -y python3.6 curl python-pip python-dev build-essential 
	
RUN python3.6 --version
RUN pip --version

RUN pip install --upgrade pip

WORKDIR $MODELPATH 
EXPOSE 5000
`
