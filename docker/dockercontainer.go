package docker

import (
	"fmt"
	"io"
	"os"

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

	fmt.Println(json)
}
