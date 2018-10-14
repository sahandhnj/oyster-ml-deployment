package docker

import (
	"io"

	"github.com/sahandhnj/apiclient/filemanager"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

const (
	DockerFileDefaultName = "Dockerfile"
)

func (c *DockerCli) BuildImage(dockerFile string, tag []string) (io.Reader, error) {
	ctx := context.Background()

	fm, err := filemanager.NewFileStoreManager()
	if err != nil {
		return nil, err
	}

	dockerFileTarReader, err := fm.TarFile(dockerFile)
	if err != nil {
		return nil, err
	}

	imageBuildRes, err := c.cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dockerFile,
			Remove:     true,
			Tags:       tag,
		})

	if err != nil {
		return nil, err
	}

	return imageBuildRes.Body, nil
}
