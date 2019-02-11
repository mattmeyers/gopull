package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mattmeyers/gopull"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func handleConfig(c *cli.Context) error {
	reposDir := c.String("repos-dir")
	gopullDir := c.String("gopull-dir")
	scriptsDir := c.String("scripts-dir")

	if reposDir != "" {
		viper.Set("repos_dir", reposDir)
		err := viper.WriteConfig()
		if err != nil {
			log.Fatalf("Failed to write repos_dir to config\nerr: %s", err)
		}
	}

	if gopullDir != "" {
		viper.Set("gopull_dir", gopullDir)
		err := viper.WriteConfig()
		if err != nil {
			log.Fatalf("Failed to write gopull_dir to config\nerr: %s", err)
		}
	}

	if scriptsDir != "" {
		viper.Set("scripts_dir", scriptsDir)
		err := viper.WriteConfig()
		if err != nil {
			log.Fatalf("Failed to write scripts_dir to config\nerr: %s", err)
		}
	}

	for _, key := range viper.AllKeys() {
		if path := strings.Split(key, "."); len(path) > 1 && path[0] == "paths" {
			fmt.Printf("%s=%s\n", path[1], viper.GetString(key))
		}
	}

	return nil
}

func handleList(c *cli.Context) error {
	repos := gopull.GetAllLocalRepos()
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

	// uri is in the form git@<remote>:<user>/<repo>. This reformats to
	// the path <remote>/<user>/<repo>
	fullName := strings.SplitAfterN(uri, "@", 2)[1]
	fullName = strings.Replace(fullName, ":", "/", 1)
	fullName = strings.Replace(fullName, ".git", "", 1)
	fullName = strings.Replace(fullName, ".org", "", 1)
	fullName = strings.Replace(fullName, ".com", "", 1)
	repoPathVars := strings.SplitN(fullName, "/", 3)

	repo := gopull.LocalRepo{
		Remote:           repoPathVars[0],
		User:             repoPathVars[1],
		Name:             repoPathVars[2],
		FullName:         fullName,
		Branch:           branch,
		Path:             fmt.Sprintf("%s/%s", viper.GetString("paths.repos_dir"), fullName),
		DeploymentScript: fmt.Sprintf("%s/deployment_scripts/%s_deploy.sh", viper.GetString("paths.gopull_dir"), repoPathVars[2]),
	}

	gopull.GitClone(uri, repo)
	repo.AddLocalRepo()
	repo.InitDeploymentScript()

	return nil
}

func handleEdit(c *cli.Context) error {
	fmt.Println("This has not been implemented yet. I only have so much time. Don't judge me.")
	return nil
}

func handleDelete(c *cli.Context) error {
	shouldDelete := c.Bool("r")

	if len(c.Args()) != 1 {
		cli.ShowCommandHelpAndExit(c, "delete", 1)
	}

	repo, err := gopull.DeleteLocalRepo(c.Args().First())
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
