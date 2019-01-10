package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/urfave/cli"
)

var _dir string

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s", err))
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
		viper.Set("repos_dir", reposDir)
		err := viper.WriteConfig()
		if err != nil {
			log.Fatalf("Failed to write to config\nerr: %s", err)
		}
	}

	for _, key := range viper.AllKeys() {
		fmt.Printf("%s=%s\n", key, viper.GetString(key))
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
		Path:             fmt.Sprintf("%s/%s/%s", viper.GetString("repos_dir"), user, name),
		DeploymentScript: fmt.Sprintf("%s/deployment_scripts/%s_deploy.sh", viper.GetString("gopull_dir"), name),
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
