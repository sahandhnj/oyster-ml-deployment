package service

import (
	"fmt"
	"path"
	"path/filepath"
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
	RunFile              = "run.sh"
	TmpServerFile        = "server.py"
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
	port := 5000

	versions, err := vs.DBHandler.VersionService.VersionsByModelId(vs.Model.ID)
	if err != nil {
		return err
	}

	for _, v := range versions {
		if v.VersionNumber > versionNumber {
			versionNumber = v.VersionNumber
		}
		if v.Port > port {
			port = v.Port
		}
	}
	versionNumber = versionNumber + 1
	port = port + 1

	version, err := types.NewVersion(versionNumber, vs.Model.ID)
	if err != nil {
		return err
	}
	version.ID = vs.DBHandler.VersionService.GetNextIdentifier()
	version.Port = port

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

	fmt.Printf("%s\t%s\t%s\t%s\n", "Name", "Version number", "Deployed", "Image Tag")
	for _, ver := range versions {
		fmt.Printf("%s\t%d\t%t\t\t%s\n", ver.Name, ver.VersionNumber, ver.Deployed, ver.ImageTag)
	}

	return nil
}

func (vs *VersionService) Deploy(versionNumber int, dcli *docker.DockerCli, verbose bool) error {
	version, err := vs.DBHandler.VersionService.VersionByVersionNumber(versionNumber, vs.Model.ID)
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

	version.ImageTag = mainTag
	containerName := vs.Model.Name + "" + strconv.Itoa(version.VersionNumber) + "-api"

	mountPath, err := filepath.Abs(vs.file.GetStorePath(version.Name))
	if err != nil {
		return err
	}

	containerId, err := dcli.CreateContainer(containerName, version.ImageTag, mountPath, strconv.Itoa(version.Port))
	if err != nil {
		return err
	}

	version.ContainerId = containerId

	if version.RedisEnabled {
		redisContainerName := vs.Model.Name + "" + strconv.Itoa(version.VersionNumber) + "-redis"

		redisContainerId, err := dcli.CreateRedisContainer(redisContainerName)
		if err != nil {
			return err
		}

		version.RedisContainerId = redisContainerId

		networkName := vs.Model.Name + "" + strconv.Itoa(version.VersionNumber) + "-network"
		networkId, err := dcli.CreateNetwork(networkName)
		if err != nil {
			return err
		}

		version.NetworkId = networkId

		dcli.ConnectToNetwork(networkId, containerId)
		dcli.ConnectToNetwork(networkId, redisContainerId)
	}

	vs.DBHandler.VersionService.UpdateVersion(version.ID, version)

	return nil
}

func (vs *VersionService) Start(versionNumber int, dcli *docker.DockerCli) error {
	version, err := vs.DBHandler.VersionService.VersionByVersionNumber(versionNumber, vs.Model.ID)
	if err != nil {
		return err
	}

	if version.RedisEnabled {
		dcli.ContainerStart(version.RedisContainerId)
	}

	dcli.ContainerStart(version.ContainerId)

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
	fm.CopyToStore(path.Join(vs.Model.ModelPath, TmpServerFile), path.Join(vs.Version.Name, TmpServerFile))
	fm.CopyToStore(path.Join(vs.Model.ModelPath, RunFile), path.Join(vs.Version.Name, RunFile))

	vs.createDockerFile()

	return nil
}

func (vs *VersionService) createDockerFile() {
	docker_file_static = docker_file_static + "EXPOSE " + strconv.Itoa(vs.Version.Port) + "\n"
	docker_file_static = docker_file_static + "RUN pip install --user " + vs.file.ReadRQLineByLine(path.Join(vs.Version.Name, RequirementsFile))
	docker_file_static = docker_file_static + "\nCMD bash run.sh"

	vs.file.WriteToFile(path.Join(vs.Version.Name, DockerFileName), docker_file_static)
	vs.file.WriteToFile(path.Join(vs.Version.Name, DockerIgnoreFileName), "")
}

var docker_file_static = `FROM ubuntu:18.04
ENV MODELPATH /src

RUN apt-get update && apt-get install -y software-properties-common

RUN add-apt-repository ppa:jonathonf/python-3.6
RUN apt-get update && apt-get install -y python3.6 curl python-pip python-dev build-essential

RUN pip install --upgrade pip

WORKDIR $MODELPATH
`
