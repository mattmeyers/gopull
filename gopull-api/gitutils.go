package main

import (
	"fmt"
	"log"
	"os/exec"
)

func GitPull(repo LocalRepo) {
	cmd := exec.Command(fmt.Sprintf("./deployment_scripts/%s", repo.DeploymentScript), repo.Path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.CombinedOutput() failed with %s\n", err)
	}
	fmt.Printf("%s\n", string(out))
}
