package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
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

	fmt.Printf("Repo Name: %s\n", info.Repository.Name)
	fmt.Printf("Branch Name: %s\n", info.Push.Changes[0].New.Name)

	GitPull("/home/matt/webhook_test")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(""); err != nil {
		panic(err)
	}
}
