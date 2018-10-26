package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sahandhnj/apiclient/backend"

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
			Name:    "truncate",
			Aliases: []string{"tr"},
			Usage:   "delete all versions of a model and the model itself",
			Action: func(c *cli.Context) error {
				modelservice, err := service.NewModelService(nil, dbhandler)
				if err != nil {
					log.Fatal(err)
				}

				versionService, err := service.NewVersionService(modelservice.Model, dbhandler)
				if err != nil {
					log.Fatal(err)
				}
				modelservice.VersionService = versionService

				dc := docker.NewDockerCli()

				err = modelservice.Truncate(modelservice.Model.ID, dc)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
		},
		{
			Name:    "dev",
			Aliases: []string{"d"},
			Usage:   "development control",
			Subcommands: []cli.Command{
				{
					Name:  "server",
					Usage: "check status of version",
					Action: func(c *cli.Context) error {
						modelservice, err := service.NewModelService(nil, dbhandler)
						if err != nil {
							log.Fatal(err)
						}

						versionService, err := service.NewVersionService(modelservice.Model, dbhandler)
						if err != nil {
							log.Fatal(err)
						}

						modelservice.VersionService = versionService

						server := &backend.Server{
							Address:        ":3000",
							DbHandler:      dbhandler,
							ModelService:   modelservice,
							VersionService: versionService,
						}

						err = server.Start()
						if err != nil {
							log.Fatal(err)
						}

						return nil
					},
				},
				{
					Name:  "test",
					Usage: "test new features",
					Action: func(c *cli.Context) error {
						modelservice, err := service.NewModelService(nil, dbhandler)
						if err != nil {
							log.Fatal(err)
						}

						// versionService, err := service.NewVersionService(modelservice.Model, dbhandler)
						// if err != nil {
						// 	log.Fatal(err)
						// }

						modelservice.GetAll()

						return nil
					},
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

						dc := docker.NewDockerCli()
						err = versionService.PrintVersions(dc)
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
				{
					Name:  "stop",
					Usage: "stop containers of version",
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

						err = versionService.Stop(vnum, dc)
						if err != nil {
							log.Fatal(err)
						}

						return nil
					},
				},
				{
					Name:  "down",
					Usage: "delete containers of version",
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

						err = versionService.Down(vnum, dc, false)
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
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
