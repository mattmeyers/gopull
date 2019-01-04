package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func GitClone(uri string, repo LocalRepo) {
	mkdir := exec.Command("mkdir", fmt.Sprintf("%s/%s", os.Getenv("REPOS_DIR"), repo.User))
	out, err := mkdir.CombinedOutput()
	if err != nil {
		log.Printf("cmd.CombinedOutput() failed with %s\n", err)
	} else {
		log.Printf("Created directory: %s/%s\n", os.Getenv("REPOS_DIR"), repo.User)
	}

	clone := exec.Command("git", "clone", "--single-branch", "--branch", repo.Branch, uri)
	clone.Dir = os.Getenv("REPOS_DIR") + "/" + repo.User
	out, err = clone.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.CombinedOutput() failed with %s\nDid you remember to add your SSH key to the remote?", err)
	}
	fmt.Printf("%s\n", string(out))
}
