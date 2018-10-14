package service

import (
	"fmt"
	"path"
	"strconv"

	"github.com/sahandhnj/apiclient/docker"

	"github.com/sahandhnj/apiclient/db"
	"github.com/sahandhnj/apiclient/filemanager"
	"github.com/sahandhnj/apiclient/types"
)

const (
	RequirementsFile     = "requirements.txt"
	ModelTarFile         = "model.tar.gz"
	DockerFileName       = "DockerFile"
	DockerIgnoreFileName = ".dockerignore"
	BuildLogFile         = "buildlog"
)

type VersionService struct {
	Version   *types.Version
	Model     *types.Model
	file      *filemanager.FileStoreManager
	DBHandler *db.DBStore
}

func NewVersionService(model *types.Model, dbHandler *db.DBStore) (*VersionService, error) {
	file, err := filemanager.NewFileStoreManager()
	if err != nil {
		return nil, err
	}

	versionService := VersionService{
		file:      file,
		DBHandler: dbHandler,
		Model:     model,
	}

	return &versionService, nil
}

func (vs *VersionService) NewVersion() error {
	versionNumber := 0
	versions, err := vs.DBHandler.VersionService.VersionsByModelId(vs.Model.ID)
	if err != nil {
		return err
	}

	for _, v := range versions {
		if v.VersionNumber > versionNumber {
			versionNumber = v.VersionNumber
		}
	}
	versionNumber = versionNumber + 1

	version, err := types.NewVersion(versionNumber, vs.Model.ID)
	if err != nil {
		return err
	}
	version.ID = vs.DBHandler.VersionService.GetNextIdentifier()

	err = vs.DBHandler.VersionService.CreateVersion(version)
	if err != nil {
		return err
	}

	vs.Version = version
	err = vs.Apply()
	if err != nil {
		return err
	}

	return nil
}

func (vs *VersionService) PrintVersions() error {
	versions, err := vs.DBHandler.VersionService.VersionsByModelId(vs.Model.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%s\t%s\t%s\n", "Name", "Version number", "Deployed")
	for _, ver := range versions {
		fmt.Printf("%s\t%d\t%t\n", ver.Name, ver.VersionNumber, ver.Deployed)
	}

	return nil
}

func (vs *VersionService) Deploy(versionNumber int, dcli *docker.DockerCli, verbose bool) error {
	version, err := vs.DBHandler.VersionService.VersionByVersionNumber(versionNumber)
	if err != nil {
		return err
	}

	dockerFilePath := path.Join(vs.file.GetStorePath(version.Name), DockerFileName)
	mainTag := "oyster/" + vs.Model.Name + ":" + strconv.Itoa(version.VersionNumber)
	tags := []string{mainTag}

	fmt.Println("Deploying: ")
	fmt.Println(tags)

	logs, err := dcli.BuildImage(dockerFilePath, tags)
	if err != nil {
		return err
	}

	logFilePath := path.Join(vs.file.GetStorePath(version.Name), BuildLogFile)

	fmt.Println("Writing image build logs into: ", logFilePath)
	err = vs.file.WriteToFileWithReader(logFilePath, logs)
	if err != nil {
		return err
	}

	if verbose {
		err = vs.file.StreamFileToStdOut(logFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vs *VersionService) Apply() error {
	fm, err := filemanager.NewFileStoreManager()
	if err != nil {
		return err
	}

	fm.CreateDirectoryInStore(vs.Version.Name)
	fm.CTarGz(path.Join(vs.Version.Name, ModelTarFile), []string{vs.Model.ModelPath}, false)
	fm.CopyToStore(path.Join(vs.Model.ModelPath, RequirementsFile), path.Join(vs.Version.Name, RequirementsFile))
	vs.createDockerFile()

	return nil
}

func (vs *VersionService) createDockerFile() {
	docker_file_static = docker_file_static + "RUN pip install --user " + vs.file.ReadRQLineByLine(path.Join(vs.Version.Name, RequirementsFile))
	vs.file.WriteToFile(path.Join(vs.Version.Name, DockerFileName), docker_file_static)
	vs.file.WriteToFile(path.Join(vs.Version.Name, DockerIgnoreFileName), "")
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
