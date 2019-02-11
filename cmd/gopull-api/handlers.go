package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/mattmeyers/gopull"
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
	repoFullName := info.Respository.FullName
	branchName := strings.SplitN(info.Ref, "/", 3)[2]

	fmt.Printf("Github Repo Full Name: %s\n", repoFullName)
	fmt.Printf("Github Branch Name: %s\n", branchName)

	localRepo := gopull.GetLocalRepo(repoFullName)

	if branchName == localRepo.Branch {
		gopull.GitPull(localRepo)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(fmt.Sprintf("Successfuly pulled %s on branch %s", repoName, branchName)); err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(fmt.Sprintf("Code pushed to branch %s, local repo is on branch %s", branchName, localRepo.Branch)); err != nil {
			panic(err)
		}
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
	repoFullName := info.Repository.FullName
	branchName := info.Push.Changes[0].New.BranchName

	fmt.Printf("Bitbucket Repo Full Name: %s\n", repoFullName)
	fmt.Printf("Bitbucket Branch Name: %s\n", branchName)

	localRepo := gopull.GetLocalRepo(repoFullName)

	if branchName == localRepo.Branch {
		gopull.GitPull(localRepo)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(fmt.Sprintf("Successfuly pulled %s on branch %s", repoName, branchName)); err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(fmt.Sprintf("Code pushed to branch %s, local repo is on branch %s", branchName, localRepo.Branch)); err != nil {
			panic(err)
		}
	}
}

func ReceiveGitlab(w http.ResponseWriter, r *http.Request) {
	var info GitlabWebhook

	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	repoName := info.Project.Name
	repoFullName := info.Project.FullName
	branchName := strings.SplitN(info.Ref, "/", 3)[2]

	fmt.Printf("Gitlab Repo Full Name: %s\n", repoFullName)
	fmt.Printf("Gitlab Branch Name: %s\n", branchName)

	localRepo := gopull.GetLocalRepo(repoFullName)

	if branchName == localRepo.Branch {
		gopull.GitPull(localRepo)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(fmt.Sprintf("Successfuly pulled %s on branch %s", repoName, branchName)); err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(fmt.Sprintf("Code pushed to branch %s, local repo is on branch %s", branchName, localRepo.Branch)); err != nil {
			panic(err)
		}
	}
}
