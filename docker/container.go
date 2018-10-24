package docker

import (
	"fmt"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

type ContainerState int

const (
	Created    ContainerState = 0
	Running    ContainerState = 1
	Paused     ContainerState = 2
	Restarting ContainerState = 3
	Removing   ContainerState = 4
	Exited     ContainerState = 5
	Dead       ContainerState = 6
	Unknown    ContainerState = 7
)

func (st ContainerState) String() string {
	states := [...]string{
		"created",
		"running",
		"paused",
		"restarting",
		"removing",
		"exited",
		"dead"}

	if st < Created || st > Dead {
		return "Unknown"
	}

	return states[st]
}

func getState(str string) ContainerState {
	switch str {
	case "created":
		return Created
	case "running":
		return Running
	case "paused":
		return Paused
	case "restarting":
		return Restarting
	case "removing":
		return Removing
	case "exited":
		return Exited
	case "dead":
		return Dead
	default:
		return Unknown
	}
}

func (c *DockerCli) CreateContainer(name string, imageTag string, mountPath string, port string) (string, error) {
	ctx := context.Background()
	redisIP := "generic_redis"
	resp, err := c.cli.ContainerCreate(ctx, &container.Config{
		Hostname: name,
		// Domainname:   hostName,
		Image:        imageTag,
		Tty:          true,
		AttachStderr: true,
		AttachStdout: true,
		ExposedPorts: nat.PortSet{
			nat.Port("5000"): {},
		},
		Env: []string{
			fmt.Sprintf("REDIS_HOST=http://%s:6379", redisIP),
		},
		// Labels: map[string]string{
		// 	"rsc": hostName,
		// },
	}, &container.HostConfig{
		AutoRemove: false,
		PortBindings: nat.PortMap{
			nat.Port("5000"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: port}},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: mountPath,
				Target: "/src",
			},
		},
	}, &network.NetworkingConfig{}, name)

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *DockerCli) CreateRedisContainer(name string) (string, error) {
	ctx := context.Background()

	resp, err := c.cli.ContainerCreate(ctx, &container.Config{
		Hostname: name,
		Image:    "redis",
	}, &container.HostConfig{
		AutoRemove: false,
	}, &network.NetworkingConfig{}, name)

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *DockerCli) ContainerStart(id string) error {
	ctx := context.Background()

	err := c.cli.ContainerStart(ctx, id, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *DockerCli) ContainerStop(id string) error {
	ctx := context.Background()

	err := c.cli.ContainerStop(ctx, id, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *DockerCli) ContainerDelete(id string) error {
	ctx := context.Background()

	err := c.cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	return nil
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

func (c *DockerCli) IsContainerRunning(id string) (string, error) {
	containers, err := c.cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})

	if err != nil {
		return "", err
	}

	for _, container := range containers {
		if id == container.ID {
			return getState(container.State).String(), nil
		}
	}

	return "", nil
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
