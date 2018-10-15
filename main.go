package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sahandhnj/apiclient/db"
	"github.com/sahandhnj/apiclient/docker"
	"github.com/sahandhnj/apiclient/service"
	"github.com/sahandhnj/apiclient/types"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	dbhandler, err := db.NewDBStore()
	if err != nil {
		log.Fatal(err)
	}

	CurrentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

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
				model := types.NewModel(name, description, modelPath, CurrentDir)

				modelservice, err := service.NewModelService(model, dbhandler)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Project with following settings has been initialized")
				modelservice.Model.PrintInfo()

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
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "version control",
			Subcommands: []cli.Command{
				{
					Name:  "commit",
					Usage: "commit a new version",
					Action: func(c *cli.Context) error {
						modelservice, err := service.NewModelService(nil, dbhandler)
						if err != nil {
							log.Fatal(err)
						}

						versionService, err := service.NewVersionService(modelservice.Model, dbhandler)
						if err != nil {
							log.Fatal(err)
						}

						err = versionService.NewVersion()
						if err != nil {
							log.Fatal(err)
						}

						fmt.Printf("Version %d has been commited into %s", versionService.Version.VersionNumber, versionService.Version.Name)
						if err != nil {
							log.Fatal(err)
						}

						return nil
					},
				},
				{
					Name:  "list",
					Usage: "list versions",
					Action: func(c *cli.Context) error {
						modelservice, err := service.NewModelService(nil, dbhandler)
						if err != nil {
							log.Fatal(err)
						}

						versionService, err := service.NewVersionService(modelservice.Model, dbhandler)
						if err != nil {
							log.Fatal(err)
						}

						err = versionService.PrintVersions()
						if err != nil {
							log.Fatal(err)
						}

						return nil
					},
				},
				{
					Name:  "deploy",
					Usage: "deploy version",
					Action: func(c *cli.Context) error {
						versionNumber := c.Args().Get(0)
						if versionNumber == "" {
							log.Fatal("You have to specify the version number")
						}

						modelservice, err := service.NewModelService(nil, dbhandler)
						if err != nil {
							log.Fatal(err)
						}

						versionService, err := service.NewVersionService(modelservice.Model, dbhandler)
						if err != nil {
							log.Fatal(err)
						}

						dc := docker.NewDockerCli()

						vnum, err := strconv.Atoi(versionNumber)
						if err != nil {
							log.Fatal(err)
						}

						err = versionService.Deploy(vnum, dc, c.Bool("verbose"))
						if err != nil {
							log.Fatal(err)
						}

						return nil
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "verbose",
							Value: "false",
						},
					},
				},
				{
					Name:  "start",
					Usage: "start containers of version",
					Action: func(c *cli.Context) error {
						versionNumber := c.Args().Get(0)
						if versionNumber == "" {
							log.Fatal("You have to specify the version number")
						}

						modelservice, err := service.NewModelService(nil, dbhandler)
						if err != nil {
							log.Fatal(err)
						}

						versionService, err := service.NewVersionService(modelservice.Model, dbhandler)
						if err != nil {
							log.Fatal(err)
						}

						dc := docker.NewDockerCli()

						vnum, err := strconv.Atoi(versionNumber)
						if err != nil {
							log.Fatal(err)
						}

						err = versionService.Start(vnum, dc)
						if err != nil {
							log.Fatal(err)
						}

						return nil
					},
				},
			},
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "show config",
			Action: func(c *cli.Context) error {
				modelservice, err := service.NewModelService(nil, dbhandler)
				if err != nil {
					log.Fatal(err)
				}

				modelservice.Model.PrintInfo()

				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test",
			Action: func(c *cli.Context) error {
				// dbhandler := db.NewDBStore()
				// model, err := model.NewModel("name", "description", "model")
				// if err != nil {
				// 	log.Fatal(err)
				// }

				// model.PrintInfo()

				return nil
			},
		},
		{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "deploy the model",
			Action: func(c *cli.Context) error {
				// model, err := model.ReadModel()
				// if err != nil {
				// 	log.Fatal(err)
				// }

				// model.PrintInfo()

				// fmt.Println("Deploying docker image")

				// dc := docker.NewDockerCli()
				// names := dc.DeployStack(model.Config.Name)
				// nodes := make([]*node.Node, 0)

				// for _, id := range names {
				// 	fmt.Println("Reading logs of: ", dc.GetContainerImageName(id))

				// 	_, node := node.NewNode(dc.GetContainerImageName(id), model.Config.ID, id, node.Running)
				// 	nodes = append(nodes, node)
				// 	go dc.ShowLogs(id)
				// }

				// model.Nodes = nodes

				// fmt.Print("Press 'Enter' to continue...")
				// bufio.NewReader(os.Stdin).ReadBytes('\n')
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

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
