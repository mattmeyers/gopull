package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func handleConfig(c *cli.Context) error {
	reposDir := c.String("repos-dir")
	gopullDir := c.String("gopull-dir")

	if reposDir != "" {
		viper.Set("repos_dir", reposDir)
		err := viper.WriteConfig()
		if err != nil {
			log.Fatalf("Failed to write to config\nerr: %s", err)
		}
	}

	if gopullDir != "" {
		viper.Set("gopull_dir", gopullDir)
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
	repo.AddLocalRepo()
	repo.InitDeploymentScript()

	return nil
}

func handleEdit(c *cli.Context) error {
	fmt.Println("Edited repo")
	return nil
}

func handleDelete(c *cli.Context) error {
	shouldDelete := c.Bool("r")

	if len(c.Args()) != 1 {
		cli.ShowCommandHelpAndExit(c, "delete", 1)
	}

	repo, err := DeleteLocalRepo(c.Args().First())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if err = os.Remove(os.ExpandEnv(repo.DeploymentScript)); err != nil {
		panic(fmt.Errorf("fatal error deleting deployment script: %s", err))
	}

	if shouldDelete {
		err = os.RemoveAll(os.ExpandEnv(repo.Path))
		if err != nil {
			panic(fmt.Errorf("fatal error deleting repo directory: %s", err))
		}
		fmt.Println("The repository directory has been deleted. You may still have to do additional cleanup such as removing process manager configurations.")
	}

	fmt.Printf("GoPull is no longer managing %s.\n", repo.Name)

	return nil
}
