package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/sahandhnj/apiclient/docker"
	"github.com/sahandhnj/apiclient/pearl"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "Oysterbox"
	app.Usage = "We deploy everyting"
	app.Version = "0.0.1"

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
				model := "model"
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

				return nil
			},
		},
		{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "deploy the model",
			Action: func(c *cli.Context) error {
				pearl, _ := pearl.NewPearl("prj1", "some AI model", "model.hdf5")
				pearl.PrintInfo()
				pearl.Config()

				fmt.Println("Deploying docker image")

				dc := docker.NewDockerCli()
				names := dc.DeployStack()

				for _, id := range names {
					fmt.Println("Reading logs of: ", id)
					go dc.ShowLogs(id)
				}

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

				dc := docker.NewDockerCli()
				dc.ListContainers()

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
