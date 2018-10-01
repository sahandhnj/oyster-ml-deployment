package docker

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type DockerCli struct {
	cli *client.Client
}

func NewDockerCli() *DockerCli {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))

	if err != nil {
		panic(err)
	}

	return &DockerCli{cli}
}

func (c *DockerCli) ListContainers() {
	containers, err := c.cli.ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.Names)
	}
}

func (c *DockerCli) ListImages() {
	images, err := c.cli.ImageList(context.Background(), types.ImageListOptions{})

	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Println(image.Labels)
	}
}

func (c *DockerCli) CreateContainer() {
	ctx := context.Background()

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFile := "Dockerfile"
	dockerFileReader, err := os.Open("Dockerfile")
	if err != nil {
		log.Fatal(err, " :unable to open Dockerfile")
	}

	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		log.Fatal(err, " :unable to read dockerfile")
	}

	tarHeader := &tar.Header{
		Name: dockerFile,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		log.Fatal(err, " :unable to write tar header")
	}

	_, err = tw.Write(readDockerFile)
	if err != nil {
		log.Fatal(err, " :unable to write tar body")
	}

	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	imageBuildRes, err := c.cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dockerFile,
			Remove:     true,
			Tags:       []string{"oysterimage"},
		})
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}

	defer imageBuildRes.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildRes.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}

	resp, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image: "my-docker-image",
		Cmd:   []string{"python", "file.py"},
		Tty:   true,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := c.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := c.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := c.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
