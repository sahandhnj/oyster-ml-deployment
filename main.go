package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/sahandhnj/apiclient/types/model"
	"github.com/sahandhnj/apiclient/types/node"

	"github.com/sahandhnj/apiclient/docker"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "Oysterbox"
	app.Usage = "We deploy everyting"
	app.Version = "0.0.0.3"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "gpu, g",
			Value: "gpu support",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"a"},
			Usage:   "start a project",
			Action: func(c *cli.Context) error {
				name := c.String("name")
				modelPath := c.String("modelPath")
				description := c.String("description")

				model, err := model.NewModel(name, description, modelPath)
				if err != nil {
					fmt.Print(err)
				}

				model.PrintInfo()

				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Value: "pearl",
				},
				cli.StringFlag{
					Name:  "modelPath",
					Value: "model",
				},
				cli.StringFlag{
					Name:  "description",
					Value: "This is a test model",
				},
			},
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "show config",
			Action: func(c *cli.Context) error {
				model, err := model.ReadModel()
				if err != nil {
					fmt.Print(err)
				}

				model.PrintInfo()

				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test",
			Action: func(c *cli.Context) error {
				dc := docker.NewDockerCli()
				dc.GetContainerConfig("3667bb131b30")
				return nil
			},
		},
		{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "deploy the model",
			Action: func(c *cli.Context) error {
				model, err := model.ReadModel()
				if err != nil {
					fmt.Print(err)
				}

				model.PrintInfo()

				fmt.Println("Deploying docker image")

				dc := docker.NewDockerCli()
				names := dc.DeployStack(model.Config.Name)
				nodes := make([]*node.Node, 0)

				for _, id := range names {
					fmt.Println("Reading logs of: ", dc.GetContainerImageName(id))

					_, node := node.NewNode(dc.GetContainerImageName(id), model.Config.ID, id, node.Running)
					nodes = append(nodes, node)
					go dc.ShowLogs(id)
				}

				model.Nodes = nodes

				fmt.Print("Press 'Enter' to continue...")
				bufio.NewReader(os.Stdin).ReadBytes('\n')
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list APIs",
			Action: func(c *cli.Context) error {
				fmt.Println("List of docker containers")

				// TODO

				return nil
			},
		},
		{
			Name:    "drop",
			Aliases: []string{"l"},
			Usage:   "drop all container APIs",
			Action: func(c *cli.Context) error {
				fmt.Println("Dropping docker containers")

				// TODO

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
