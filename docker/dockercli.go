package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (c *DockerCli) PrintAllImages() {
	images, err := c.cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("LIST IMAGES\n-----------------------")
	fmt.Println("Image ID | Repo Tags | Size")
	for _, image := range images {
		fmt.Printf("%s | %s | %d\n", image.ID, image.RepoTags, image.Size)
	}

}

func (c *DockerCli) PrintAllCointainers() {
	containers, err := c.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("LIST CONTAINERS\n-----------------------")
	fmt.Println("Container Names | Image | Mounts")

	for _, container := range containers {
		fmt.Printf("%s | %s | %s\n", container.Names, container.Image, container.Mounts)
	}
}

func (c *DockerCli) PrintAllNetworks() {
	networks, err := c.cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Print("\n\n\n")
	fmt.Println("LIST NETWORKS\n-----------------------")
	fmt.Println("Network Name | ID")
	for _, network := range networks {
		fmt.Printf("%s | %s\n", network.Name, network.ID)
	}
}
