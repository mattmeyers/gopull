package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"runtime"
	"path"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var _dir string

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("No caller information")
	}
	_dir = path.Dir(filename)

	err := godotenv.Load(fmt.Sprintf("%s/.env", _dir))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := cli.NewApp()
	app.Name = "gopull"
	app.Usage = "Configure the GoPull REST API to pull remote repo changes"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:   "config",
			Usage:  "Configure the GoPull environment",
			Action: handleConfig,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "repos-dir, r",
					Usage: "Set the base directory where repositories are located. Defaults to \"$HOME/repos\"",
				},
				// cli.StringFlag{
				// 	Name:  "gopull-dir, g",
				// 	Usage: "Set the GoPull API directory. Defaults to \"$GOPATH/src/gopull\"",
				// },
			},
		},
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

func handleConfig(c *cli.Context) error {
	reposDir := c.String("repos-dir")

	if reposDir != "" {
		env := map[string]string{
			"REPOS_DIR":  reposDir,
			"GOPULL_DIR": os.Getenv("GOPULL_DIR"),
		}

		err := godotenv.Write(env, fmt.Sprintf("%s/.env", _dir))
		if err != nil {
			log.Fatalf("Could not write to .env\nerr: %s", err)
		}

		return nil
	}

	var env map[string]string
	env, err := godotenv.Read(fmt.Sprintf("%s/.env", _dir))
	if err != nil {
		log.Fatalf("Could not read .env file\nerr: %s", err)
	}

	for key, val := range env {
		fmt.Printf("%s=%s\n", key, val)
	}

	return nil
}

func handleList(c *cli.Context) error {
	repos := GetAllLocalRepos()
	for _, repo := range repos {
		fmt.Println(repo.FullName)
	}
	return nil
}

func handleAdd(c *cli.Context) error {
	uri := c.String("uri")
	branch := c.String("branch")

	if uri == "" || branch == "" {
		cli.ShowCommandHelpAndExit(c, "add", 1)
	}

	fullName := strings.Replace(strings.SplitN(uri, ":", 2)[1], ".git", "", 1)
	repoPathVars := strings.SplitN(fullName, "/", 2)
	user := repoPathVars[0]
	name := repoPathVars[1]

	repo := LocalRepo{
		User:             user,
		Name:             name,
		FullName:         fullName,
		Branch:           branch,
		Path:             fmt.Sprintf("%s/%s/%s", os.Getenv("REPOS_DIR"), user, name),
		DeploymentScript: fmt.Sprintf("%s/deployment_scripts/%s_deploy.sh", os.Getenv("GOPULL_DIR"), name),
	}

	GitClone(uri, repo)
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
