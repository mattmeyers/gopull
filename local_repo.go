package gopull

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	var repos map[string]map[string]map[string]LocalRepo

	viper.UnmarshalKey("repos", &repos)

	if repos[r.Remote] == nil {
		repos[r.Remote] = map[string]map[string]LocalRepo{}
	}

	if repos[r.Remote][r.User] == nil {
		repos[r.Remote][r.User] = map[string]LocalRepo{}
	}

	repos[r.Remote][r.User][r.Name] = r

	viper.Set("repos", repos)
	if err := viper.WriteConfig(); err != nil {
		panic(fmt.Sprintf("Error writing repo to config\n%s", err))
	}

}

// InitDeploymentScript copies the sample deployment script for use with the
// new local repository.
func (r LocalRepo) InitDeploymentScript() {
	from, err := os.Open(fmt.Sprintf("%s/deploy.src.sh", os.ExpandEnv(viper.GetString("paths.scripts_dir"))))
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	newFilename := fmt.Sprintf("%s/%s_deploy.sh", os.ExpandEnv(viper.GetString("paths.scripts_dir")), r.Name)
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

// GetAllLocalRepos gets all of the managed local repos.
func GetAllLocalRepos() map[string]LocalRepo {
	var repos map[string]LocalRepo
	viper.UnmarshalKey("repos", &repos)
	return repos
}

// GetLocalRepo gets a single managed local repo.
func GetLocalRepo(name string) LocalRepo {
	var repo LocalRepo
	viper.UnmarshalKey(fmt.Sprintf("repos.%s", name), &repo)
	return repo
}

// DeleteLocalRepo deletes a managed local repo.
func DeleteLocalRepo(fullname string) (*LocalRepo, error) {
	var repo *LocalRepo
	path := strings.SplitN(fullname, "/", 3)

	viper.UnmarshalKey(fmt.Sprintf("repos.%s.%s.%s", path[0], path[1], path[2]), &repo)
	if repo == nil {
		return repo, fmt.Errorf("no managed repo with full name %s", fullname)
	}
	delete(viper.Get(fmt.Sprintf("repos.%s.%s", path[0], path[1])).(map[string]interface{}), path[2])
	if err := viper.WriteConfig(); err != nil {
		return repo, err
	}

	return repo, nil
	// repos := readInFile()
	// repo, ok := repos[repoName]
	// if ok {
	// 	delete(repos, repoName)
	// } else {
	// 	return repo, errors.New("provided name is not a managed repository")
	// }

	// writeToFile(repos)

	// return repo, nil
}

func readInFile() map[string]LocalRepo {
	filePath := fmt.Sprintf("%s/repos.json", os.ExpandEnv(viper.GetString("paths.gopull_dir")))

	var repos map[string]LocalRepo
	configFile, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to open config.json\n%s", err))
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
