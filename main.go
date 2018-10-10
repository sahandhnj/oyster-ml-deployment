package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/sahandhnj/apiclient/db"
	"github.com/sahandhnj/apiclient/docker"
	"github.com/sahandhnj/apiclient/pearl"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "Oysterbox"
	app.Usage = "We deploy everyting"
	app.Version = "0.0.0.2"
	db.Connect()

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
				name := uuid.New().String()
				model := "mlpipeline"
				description := ""

				if len(c.String("name")) > 0 {
					name = c.String("name")
				}

				if len(c.String("modelPath")) > 0 {
					model = c.String("modelPath")
				}

				if len(c.String("description")) > 0 {
					description = c.String("description")
				}

				pearl, _ := pearl.NewPearl(name, description, model)
				db.Insert(pearl)

				pearl.PrintInfo()
				pearl.Config()

				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Value: "Name of the project",
				},
				cli.StringFlag{
					Name:  "modelPath",
					Value: "path to hdf5 model",
				},
				cli.StringFlag{
					Name:  "description",
					Value: "Description of the project",
				},
			},
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "show config",
			Action: func(c *cli.Context) error {
				p, _ := pearl.ReadPearl()
				p.PrintInfo()

				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test",
			Action: func(c *cli.Context) error {
				dc := docker.NewDockerCli()
				dc.InspectContainer("ce7134284c1ebe05605e1f26b5d05cdb7a3d66d57c4bf61c185dadf2fc7cac47")
				return nil
			},
		},
		{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "deploy the model",
			Action: func(c *cli.Context) error {
				p, _ := pearl.ReadPearl()
				p.PrintInfo()

				fmt.Println("Deploying docker image")

				dc := docker.NewDockerCli()
				names := dc.DeployStack(p.Name)
				nodes := make([]*pearl.Node, 0)
				for _, id := range names {
					fmt.Println("Reading logs of: ", dc.GetContainerImageName(id))

					_, node := pearl.NewNode(dc.GetContainerImageName(id), p.ID, id, pearl.Running)
					nodes = append(nodes, node)
					go dc.ShowLogs(id)
				}

				p.Nodes = nodes
				p.UpdateConfig()

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
				pearl, _ := pearl.ReadPearl()

				dc := docker.NewDockerCli()
				dc.ListStackContainers(pearl.Name)

				return nil
			},
		},
		{
			Name:    "drop",
			Aliases: []string{"l"},
			Usage:   "drop all container APIs",
			Action: func(c *cli.Context) error {
				fmt.Println("Dropping docker containers")
				pearl, _ := pearl.ReadPearl()

				dc := docker.NewDockerCli()
				dc.Drop(pearl.Name)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
