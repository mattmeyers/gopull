package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	InitConfig()

	app := cli.NewApp()
	app.Name = "gopull-cli"
	app.Usage = "Configure the GoPull REST API to pull remote repo changes"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:   "config",
			Usage:  "Configure the GoPull environment",
			Action: handleConfig,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "repos-dir",
					Usage: "Set the base directory where repositories are located. Defaults to \"$HOME/repos\"",
				},
				cli.StringFlag{
					Name:  "gopull-dir",
					Usage: "Set the GoPull API directory. Defaults to \"$GOPATH/src/github.com/mattmeyers/gopull\"",
				},
				cli.StringFlag{
					Name:  "scripts-dir",
					Usage: "Set the deployment scripts directory. Defaults to \"$GOPATH/src/github.com/mattmeyers/gopull/deployment_scripts\"",
				},
			},
		},
		{
			Name:   "list",
			Usage:  "List configured local repos",
			Action: handleList,
		},
		{
			Name:  "add",
			Usage: "Add a new repository",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uri",
					Usage: "SSH URI of the repo. Of the form \"git@<remote>:<user>/<repository>\". (Required)",
				},
				cli.StringFlag{
					Name:  "branch",
					Usage: "Branch to be pulled. (Required)",
				},
			},
			Action: handleAdd,
		},
		{
			Name:   "edit",
			Usage:  "Edit an existing repository configuration.",
			Action: handleEdit,
		},
		{
			Name:   "delete",
			Usage:  "Delete a repository configuration.",
			Action: handleDelete,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "r",
					Usage: "Completely remove the repository. The configured directory will be completely removed.",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
