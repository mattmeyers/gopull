package gopull

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	Remote           string `json:"remote"`
	Branch           string `json:"branch"`
	Path             string `json:"path"`
	DeploymentScript string `json:"deploymentScript"`
}

// AddLocalRepo adds a new local repository configuration to the repos.json
// configuration file.
func (r LocalRepo) AddLocalRepo() {
	repos := readInFile()
	repos[r.FullName] = r
	writeToFile(repos)
}

// InitDeploymentScript copies the sample deployment script for use with the
// new local repository.
func (r LocalRepo) InitDeploymentScript() {
	from, err := os.Open(fmt.Sprintf("%s/deploy.src.sh", os.ExpandEnv(viper.GetString("scripts_dir"))))
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	newFilename := fmt.Sprintf("%s/%s_deploy.sh", os.ExpandEnv(viper.GetString("scripts_dir")), r.Name)
	to, err := os.OpenFile(newFilename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
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

// DeleteLocalRepo deletes a local repository configuration from the
// repos.json configuration file.
func DeleteLocalRepo(repoName string) (LocalRepo, error) {
	repos := readInFile()
	repo, ok := repos[repoName]
	if ok {
		delete(repos, repoName)
	} else {
		return repo, errors.New("provided name is not a managed repository")
	}

	writeToFile(repos)

	return repo, nil
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

	reposJSON, err := json.Marshal(repos)
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(filePath, reposJSON, 0644); err != nil {
		panic(err)
	}

}
