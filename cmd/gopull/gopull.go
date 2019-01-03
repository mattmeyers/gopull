package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := cli.NewApp()
	app.Name = "gopull"
	app.Usage = "Configure the GoPull REST API to pull remote repo changes"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "repo",
			Usage: "Configure local repositories",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "list, l",
					Usage: "List all currently configured local repos.",
				},
			},
			Action: handleRepo,
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "Add a new repository",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "Full name of the new repo. Of the form \"<User>/<Repo>\". (Required)",
						},
						cli.StringFlag{
							Name:  "branch",
							Usage: "Branch to be pulled. (Required)",
						},
					},
					Action: handleRepoAdd,
				},
				{
					Name:   "edit",
					Usage:  "Edit an existing repository",
					Action: handleRepoEdit,
				},
				{
					Name:   "delete",
					Usage:  "Delete a repository",
					Action: handleRepoDelete,
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func handleRepo(c *cli.Context) error {
	if c.Bool("list") {
		repos := GetAllLocalRepos()
		for _, repo := range repos {
			fmt.Println(repo.FullName)
		}
		return nil
	}

	cli.ShowSubcommandHelp(c)
	return nil
}

func handleRepoAdd(c *cli.Context) error {
	if c.String("name") == "" || c.String("branch") == "" {
		cli.ShowCommandHelpAndExit(c, "add", 1)
	}
	fmt.Printf("Repo name: %s\nBranch name: %s\n", c.String("name"), c.String("branch"))
	return nil
}

func handleRepoEdit(c *cli.Context) error {
	fmt.Println("Edited repo")
	return nil
}

func handleRepoDelete(c *cli.Context) error {
	fmt.Println("Deleted repo")
	return nil
}
