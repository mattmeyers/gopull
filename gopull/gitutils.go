package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

// GitClone creates a proper directory structure before cloning a remote repository.
//
// Remote repositories are assumed to take the form "git@<remote>:<user>/<repository>"
// and so the directory structure $REPOS_DIR/<user> is created. Then the repository is
// cloned creating the final structure of $REPOS_DIR/<user>/<repository>.
//
// Before a repository can be cloned, an SSH key must be added to the remote. Refer to
// the remote's documentation for further information:
//		- Bitbucket: https://confluence.atlassian.com/bitbucket/access-keys-294486051.html
//		- Github: https://help.github.com/articles/adding-a-new-ssh-key-to-your-github-account/
//		- Gitlab: https://docs.gitlab.com/ee/ssh/
func GitClone(uri string, repo LocalRepo) {
	mkdir := exec.Command("mkdir", fmt.Sprintf("%s/%s", os.ExpandEnv(viper.GetString("repos_dir")), repo.User))
	out, err := mkdir.CombinedOutput()
	if err != nil {
		log.Printf("cmd.CombinedOutput() failed with %s\n", err)
	} else {
		log.Printf("Created directory: %s/%s\n", os.ExpandEnv(viper.GetString("repos_dir")), repo.User)
	}

	clone := exec.Command("git", "clone", "--single-branch", "--branch", repo.Branch, uri)
	clone.Dir = os.ExpandEnv(viper.GetString("repos_dir")) + "/" + repo.User
	out, err = clone.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.CombinedOutput() failed with %s\nDid you remember to add your SSH key to the remote?", err)
	}
	fmt.Printf("%s\n", string(out))
}
