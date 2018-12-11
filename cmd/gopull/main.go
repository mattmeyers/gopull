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
	deleteRepoCommand := flag.NewFlagSet("delete-repo", flag.ExitOnError)

	// Add-Repo subcommand flag pointers
	addRepoFullNamePtr := addRepoCommand.String("fullname", "", "Full name of the remote repository. Of the form \"<User>/<Repository Name>\". (Required)")
	addRepoBranchPtr := addRepoCommand.String("branch", "", "Branch to use. (Required)")

	// Edit-Repo subcommand flag pointers
	editRepoFullNamePtr := editRepoCommand.String("fullname", "", "Full name of the remote repository. Of the form \"<User>/<Repository Name>\".")
	editRepoBranchPtr := editRepoCommand.String("branch", "", "Branch to use.")

	// Delete-Repo subcommand flag pointers
	deleteRepoFullNamePtr := deleteRepoCommand.String("fullname", "", "Full name of the remote repository. Of the form \"<User>/<Repository Name>\". (Required)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add-repo":
		addRepoCommand.Parse(os.Args[2:])
	case "edit-repo":
		editRepoCommand.Parse(os.Args[2:])
	case "delete-repo":
		deleteRepoCommand.Parse(os.Args[2:])
	default:
		printUsage()
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
	} else if deleteRepoCommand.Parsed() {
		// Required flags
		if *deleteRepoFullNamePtr == "" {
			deleteRepoCommand.PrintDefaults()
			os.Exit(1)
		}

		DeleteLocalRepo(*deleteRepoFullNamePtr)
	}

}

func printUsage() {
	fmt.Print("gopull command line tools\n\n")
	fmt.Print("Usage: gopull <command> [options]\n\n")
	fmt.Println("Available sub-commands:")
	fmt.Println("  add-repo\tAdd a new local repo configuration.")
	fmt.Println("  edit-repo\tEdit a local repo configuration.")
	fmt.Println("  delete-repo\tEdit a local repo configuration.")
}
