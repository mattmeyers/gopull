package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Bool("h", false, "Display gopull usage")
	// Subcommands
	addRepoCommand := flag.NewFlagSet("add-repo", flag.ExitOnError)
	editRepoCommand := flag.NewFlagSet("edit-repo", flag.ExitOnError)

	// Add-Repo subcommand flag pointers
	addRepoFullNamePtr := addRepoCommand.String("fullname", "", "Full name of the remote repository. Of the form \"<User>/<Repository Name>\". (Required)")
	addRepoBranchPtr := addRepoCommand.String("branch", "", "Branch to use. (Required)")

	// Edit-Repo subcommand flag pointers
	editRepoFullNamePtr := editRepoCommand.String("fullname", "", "Full name of the remote repository. Of the form \"<User>/<Repository Name>\".")
	editRepoBranchPtr := editRepoCommand.String("branch", "", "Branch to use.")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("add-repo or edit-repo subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add-repo":
		addRepoCommand.Parse(os.Args[2:])
	case "edit-repo":
		editRepoCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if addRepoCommand.Parsed() {
		// Required flags
		if *addRepoFullNamePtr == "" || *addRepoBranchPtr == "" {
			addRepoCommand.PrintDefaults()
			os.Exit(1)
		}

		fmt.Printf("Full Name: %s\nBranch: %s\n", *addRepoFullNamePtr, *addRepoBranchPtr)

		name := strings.SplitN(*addRepoFullNamePtr, "/", 2)[1]

		AddLocalRepo(LocalRepo{Name: name, FullName: *addRepoFullNamePtr, Branch: *addRepoBranchPtr, Path: fmt.Sprintf("/home/matt/%s", *addRepoFullNamePtr), DeploymentScript: fmt.Sprintf("./deployment_scripts/%s_deploy.sh", name)})
	} else if editRepoCommand.Parsed() {
		// Required flags
		if *editRepoFullNamePtr == "" && *editRepoBranchPtr == "" {
			editRepoCommand.PrintDefaults()
			os.Exit(1)
		}
	}

}
