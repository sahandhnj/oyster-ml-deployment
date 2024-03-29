package docker

import (
	"fmt"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (c *DockerCli) PrntLogs(id string) {
	ctx := context.Background()

	options := types.ContainerLogsOptions{ShowStdout: true}
	out, err := c.cli.ContainerLogs(ctx, id, options)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}

func (c *DockerCli) InspectContainer(id string) {
	ctx := context.Background()

	json, err := c.cli.ContainerInspect(ctx, id)
	if err != nil {
		panic(err)
	}

	fmt.Println(json.Config)
}

func (c *DockerCli) GetContainerImageName(id string) string {
	ctx := context.Background()

	json, err := c.cli.ContainerInspect(ctx, id)
	if err != nil {
		panic(err)
	}

	return json.Config.Image
}

func (c *DockerCli) GetContainerConfig(id string) {
	ctx := context.Background()

	json, err := c.cli.ContainerInspect(ctx, id)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", json.Config)
	// fmt.Printf("%+v\n", json.NetworkSettings)
	spew.Dump(json.Config, json.NetworkSettings)

}
