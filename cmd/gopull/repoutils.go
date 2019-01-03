package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type LocalRepo struct {
	Name             string `json:"name"`
	FullName         string `json:"fullName"`
	Branch           string `json:"branch"`
	Path             string `json:"path"`
	DeploymentScript string `json:"deploymentScript"`
}

func GetAllLocalRepos() map[string]LocalRepo {
	return readInFile()
}

func GetLocalRepo(name string) LocalRepo {
	repos := readInFile()
	return repos[name]
}

func AddLocalRepo(repo LocalRepo) {
	repos := readInFile()
	repos[repo.FullName] = repo
	writeToFile(repos)
}

func DeleteLocalRepo(repo string) {
	repos := readInFile()
	delete(repos, repo)
	writeToFile(repos)
}

func readInFile() map[string]LocalRepo {
	filePath := fmt.Sprintf("%s/repos.json", os.Getenv("GOPULL_DIR"))

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
	filePath := fmt.Sprintf("%s/repos.json", os.Getenv("GOPULL_DIR"))

	reposJson, err := json.Marshal(repos)
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(filePath, reposJson, 0644); err != nil {
		panic(err)
	}

}
