package docker

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (c *DockerCli) CreateNetwork(networkName string) (string, error) {
	ctx := context.Background()
	netresp, err := c.cli.NetworkCreate(ctx, networkName, types.NetworkCreate{})
	if err != nil {
		return "", err
	}

	return netresp.ID, nil
}

func (c *DockerCli) ConnectToNetwork(networkId string, containerId string) error {
	ctx := context.Background()

	err := c.cli.NetworkConnect(ctx, networkId, containerId, nil)
	if err != nil {
		return err
	}

	return nil
}
