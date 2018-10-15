package docker

import (
	"fmt"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"golang.org/x/net/context"
)

func (c *DockerCli) CreateContainer(hostName string, imageTag string, cmd []string, mountPath string) (string, error) {
	ctx := context.Background()

	resp, err := c.cli.ContainerCreate(ctx, &container.Config{
		Hostname: hostName,
		// Domainname:   hostName,
		Image: imageTag,
		Cmd:   cmd,
		// Tty:          true,
		// AttachStderr: true,
		// AttachStdout: true,
		// Labels: map[string]string{
		// 	"rsc": hostName,
		// },
	}, &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: mountPath,
				Target: "/src",
			},
		},
	}, nil, hostName)

	fmt.Println(resp)

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

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
