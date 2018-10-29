package docker

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

type DockerCli struct {
	cli *client.Client
}

type DockerCliParams struct {
	Host    string
	Version string
}

func NewDockerCli(params *DockerCliParams) *DockerCli {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}

	if params != nil {
		cli, err = client.NewClientWithOpts(client.WithHost(params.Host), client.WithVersion(params.Version))
		if err != nil {
			panic(err)
		}
	}

	return &DockerCli{cli}
}

func (c *DockerCli) ListContainers() error {
	containers, err := c.cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})

	if err != nil {
		return err
	}

	for _, container := range containers {
		fmt.Println("ContainerID:" + container.ID)
	}

	return nil
}

func (c *DockerCli) ListStackContainers(prjname string) []string {
	myproject, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"docker-compose.yml"},
			ProjectName:  prjname,
		},
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	names, err := myproject.Containers(context.Background(), project.Filter{
		State: project.AnyState,
	})

	if err != nil {
		log.Fatal(err)
	}

	return names
}

func (c *DockerCli) Drop(prjname string) {
	myproject, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"docker-compose.yml"},
		},
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downing the containers of the project: " + prjname)
	err = myproject.Down(context.Background(), options.Down{})

	if err != nil {
		log.Fatal(err)
	}

}

func (c *DockerCli) DeployStack(prjname string) []string {
	myproject, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"docker-compose.yml"},
			// ProjectName:  prjname,
		},
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Building the images")
	// err = myproject.Build(context.Background(), options.Build{})

	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("Upping the images")
	err = myproject.Up(context.Background(), options.Up{})

	if err != nil {
		log.Fatal(err)
	}

	names, err := myproject.Containers(context.Background(), project.Filter{
		State: project.Running,
	})

	if err != nil {
		log.Fatal(err)
	}

	return names
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

func (c *DockerCli) ShowLogs(ID string) {
	ctx := context.Background()

	out, err := c.cli.ContainerLogs(ctx, ID, types.ContainerLogsOptions{ShowStdout: true, Follow: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}

func (c *DockerCli) CreateAndRunContainer() {
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

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image: "oysterimage",
		Cmd:   []string{"python", "/workspace/hello.py"},
		Tty:   true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: path.Join(dir, "workspace"),
				Target: "",
			},
		},
	}, nil, "")

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
