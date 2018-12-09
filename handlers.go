package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func ReceiveGithub(w http.ResponseWriter, r *http.Request) {
	var info GithubWebhook

	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	repoName := info.Respository.Name
	branchName := strings.SplitN(info.Ref, "/", 3)[2]

	fmt.Printf("Github Repo Name: %s\n", info.Respository.Name)
	fmt.Printf("Github Branch Name: %s\n", branchName)

	GitPull(fmt.Sprintf("/home/matt/%s", repoName))

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(fmt.Sprintf("Successfuly pulled %s on branch %s", repoName, branchName)); err != nil {
		panic(err)
	}
}

func ReceiveBitbucket(w http.ResponseWriter, r *http.Request) {
	var info BitbucketWebhook

	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	repoName := info.Repository.Name
	branchName := info.Push.Changes[0].New.BranchName

	fmt.Printf("Bitbucket Repo Name: %s\n", repoName)
	fmt.Printf("Bitbucket Branch Name: %s\n", branchName)

	GitPull(fmt.Sprintf("/home/matt/%s", repoName))

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(fmt.Sprintf("Successfuly pulled %s on branch %s", repoName, branchName)); err != nil {
		panic(err)
	}
}
