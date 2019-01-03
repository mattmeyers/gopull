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
			Name:   "list",
			Usage:  "List configure local repos",
			Action: handleList,
		},
		{
			Name:  "add",
			Usage: "Add a new repository",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "user",
					Usage: "Owner of the repo. (Required)",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "Name of the repo. (Required)",
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
			Usage:  "Edit an existing repository",
			Action: handleEdit,
		},
		{
			Name:   "delete",
			Usage:  "Delete a repository",
			Action: handleDelete,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func handleList(c *cli.Context) error {
	repos := GetAllLocalRepos()
	for _, repo := range repos {
		fmt.Println(repo.FullName)
	}
	return nil
}

func handleAdd(c *cli.Context) error {
	user := c.String("user")
	name := c.String("name")
	branch := c.String("branch")

	if user == "" || name == "" || branch == "" {
		cli.ShowCommandHelpAndExit(c, "add", 1)
	}

	repo := LocalRepo{
		Name:             name,
		FullName:         fmt.Sprintf("%s/%s", user, name),
		Branch:           branch,
		Path:             fmt.Sprintf("%s/%s/%s", os.Getenv("REPOS_DIR"), user, name),
		DeploymentScript: fmt.Sprintf("%s/deployment_scripts/%s_deploy.sh", os.Getenv("GOPULL_DIR"), name),
	}

	AddLocalRepo(repo)

	return nil
}

func handleEdit(c *cli.Context) error {
	fmt.Println("Edited repo")
	return nil
}

func handleDelete(c *cli.Context) error {
	fmt.Println("Deleted repo")
	return nil
}
