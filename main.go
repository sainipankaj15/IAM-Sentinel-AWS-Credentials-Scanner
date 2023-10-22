package main

import (
	"fmt"

	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <path_to_git_repository>")
		return
	}

	repoPath := os.Args[1]

	// Step 1: Clone the Git repository locally
	err := cloneRepository(repoPath)
	if err != nil {
		fmt.Println("Error cloning the repository:", err)
		return
	}

	// create a folder: logs
	err = os.Mkdir("logs", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder - logs:", err)
	}

	// create a file in logs/ to log the report of scanning
	f, err := os.OpenFile("logs/"+getRepoName(repoPath)+"-result.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	rs := RepoScanner{report: f, repoPath: repoPath, cv: awsValidator{}}
	// Traverse through all files in the repository to find potential AWS IAM keys
	// through all branches and commit history
	err = rs.ScanRepo()
	if err != nil {
		fmt.Println("Error scanning repository:", err)
	}

	err = f.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
		log.Fatal(err)
	}

	// Step 3: Clean up cloned repository
	err = os.RemoveAll(getRepoName(repoPath))
	if err != nil {
		fmt.Println("Error cleaning up cloned repository:", err)
	}
}
