package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

// LocalRepo holds all information about a local repository being
// managed by GoPull.
//
// This struct is used to interact with the repos.json configuration
// file located in $GOPULL_DIR. That configuration file contains a
// list of objects that map to the LocalRepo struct.
type LocalRepo struct {
	User             string `json:"user"`
	Name             string `json:"name"`
	FullName         string `json:"fullName"`
	Branch           string `json:"branch"`
	Path             string `json:"path"`
	DeploymentScript string `json:"deploymentScript"`
}

// GetAllLocalRepos gets all of the repositories from the repos.json
// configuration file.
func GetAllLocalRepos() map[string]LocalRepo {
	return readInFile()
}

// GetLocalRepo gets a single repository configuration from the repos.json
// configuration file.
func GetLocalRepo(name string) LocalRepo {
	repos := readInFile()
	return repos[name]
}

// AddLocalRepo adds a new local repository configuration to the repos.json
// configuration file.
func AddLocalRepo(repo LocalRepo) {
	repos := readInFile()
	repos[repo.FullName] = repo
	writeToFile(repos)
}

// DeleteLocalRepo deletes a local repository configuration from the
// repos.json configuration file.
func DeleteLocalRepo(repo string) {
	repos := readInFile()
	delete(repos, repo)
	writeToFile(repos)
}

func readInFile() map[string]LocalRepo {
	filePath := fmt.Sprintf("%s/repos.json", os.ExpandEnv(viper.GetString("gopull_dir")))

	var repos map[string]LocalRepo
	configFile, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to open repos.json\n%s", err))
	}
	defer configFile.Close()

	if err := json.NewDecoder(configFile).Decode(&repos); err != nil {
		panic(fmt.Sprintf("Failed to parse json\n%s", err))
	}

	return repos
}

func writeToFile(repos map[string]LocalRepo) {
	filePath := fmt.Sprintf("%s/repos.json", os.ExpandEnv(viper.GetString("gopull_dir")))

	reposJson, err := json.Marshal(repos)
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(filePath, reposJson, 0644); err != nil {
		panic(err)
	}

}
