package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type LocalRepo struct {
	User             string `json:"user"`
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
	var repos map[string]LocalRepo
	configFile, err := os.Open("repos.json")
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
	reposJson, err := json.Marshal(repos)
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile("repos.json", reposJson, 0644); err != nil {
		panic(err)
	}

}
