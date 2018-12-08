package main

import (
	"fmt"
	"log"
	"os/exec"
)

func GitPull(path string) {
	cmd := exec.Command("git", "-C", path, "pull")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("%s\n", string(out))
}
