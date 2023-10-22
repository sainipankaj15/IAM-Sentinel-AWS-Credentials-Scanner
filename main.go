package main

import (
	"fmt"

	"log"
	"os"
)

var RepoName string

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <path_to_git_repository>")
		return
	}

	repoPath := os.Args[1]

	RepoName = repoPath

	// Cloning git repo locally
	err := cloneRepository(repoPath)
	if err != nil {
		fmt.Println("Error cloning the repository:", err)
		return
	}

	// Logs folder creating
	err = os.Mkdir("logs", os.ModePerm)

	if err != nil {
		if os.IsExist(err) {
			// Folder already exists, no need to create it again
			fmt.Println("\nLogs folder is already there")
		}else {
			fmt.Println("Error creating folder - logs:", err)
		}
	}

	// Report file inside logs 
	f, err := os.OpenFile("logs/"+getRepoName(repoPath)+"-result.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	rs := RepoScanner{report: f, repoPath: repoPath, cv: awsValidator{}}

	// Traversing all the branchs and commit 
	err = rs.ScanRepo()
	if err != nil {
		fmt.Println("Error scanning repository:", err)
	}

	err = f.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
		log.Fatal(err)
	}

	// Removing the cloned repo as it is garbage for us
	err = os.RemoveAll(getRepoName(repoPath))
	if err != nil {
		fmt.Println("Error cleaning up cloned repository:", err)
	}
}
