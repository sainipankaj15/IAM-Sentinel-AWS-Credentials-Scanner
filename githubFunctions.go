package main

import (
	"errors"
	"fmt"

	"os"
	"os/exec"
	"strings"
)

type RepoScanner struct {
	cv       CredentialValidator
	repoPath string
	report   *os.File
}

// Scans a repository for any valid secrets/credentials;
// performs a scan on each branch of the repository
func (r *RepoScanner) ScanRepo() error {
	dirName := getRepoName(r.repoPath)
	branches, err := getAllBranches(dirName)

	fmt.Println("\nAll branches are in the repoistry ", branches)

	if err != nil {
		fmt.Println("Error not getting branch ", err)
	}
	for _, branch := range branches {
		err := r.scanBranch(branch, dirName)
		if err != nil {
			fmt.Println("Error scanning branch:", err)
		}
	}

	return nil

}

// Scans a branch of repository for valid secrets;
// performs a scan on each commit of the branch
func (r *RepoScanner) scanBranch(branch, dirName string) error {
	switchToRef(branch, dirName)

	msg := fmt.Sprintf("\t \nBranch: %s\n", branch)
	fmt.Println(msg)
	r.report.Write([]byte(msg))

	commits, err := getAllCommits(dirName)
	fmt.Println("Total Commits are in this branch ", commits)

	previousVerifiedCommits := AlreadyVerifiedCommit[branch]

	diffCommits := diffCommits(commits, previousVerifiedCommits)

	AlreadyVerifiedCommit[branch] = commits

	commits = diffCommits

	if err != nil {
		return nil
	}
	for _, commit := range commits {
		err = r.scanCommit(commit, dirName)
		if err != nil {
			fmt.Println("Error scanning files in commit:", err)
			return err
		}
	}

	return nil
}

// Scans a commit for valid secrets;
// performs a scan on all directories present
func (r *RepoScanner) scanCommit(commit, dirName string) error {
	err := switchToRef(commit, dirName)
	if err != nil {
		fmt.Println("Error switching to commit:", err)
		return err
	}
	msg := fmt.Sprintf("\t\tCommit: %s\n", commit)
	fmt.Println(msg)
	r.report.Write([]byte(msg))
	err = scanDir(dirName, r.report, r.cv)
	if err != nil {
		fmt.Println("Error scanning directory:", err)
	}

	// switch back to HEAD
	cmd := exec.Command("git", "-C", dirName, "switch", "-")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error switching back to HEAD:", err)
		return err
	}

	return nil
}

func getRepoName(repoPath string) string {
	slice := strings.Split(repoPath, "/")
	folderName := slice[len(slice)-1]
	return folderName
}

func getAllBranches(dirName string) ([]string, error) {

	cmd := exec.Command("git", "ls-remote", "--heads", RepoName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	// Extracting branch names from the output
	var branches []string
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) > 1 {
			branchRef := parts[1]
			branchName := strings.TrimPrefix(branchRef, "refs/heads/")
			branches = append(branches, branchName)
		}
	}
	// fmt.Println("bRANHCES ARE ", branches)
	return branches, nil
}

func switchToRef(ref string, dirName string) error {
	cmd := exec.Command("git", "-C", dirName, "checkout", ref)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error switching to ref:", ref, err)
		return err
	}

	return nil
}

func getAllCommits(dirName string) ([]string, error) {
	cmd := exec.Command("git", "-C", dirName, "rev-list", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	commits := strings.Split(strings.TrimSpace(string(output)), "\n")
	return commits, nil
}

// Cloning Repo locally
func cloneRepository(repoPath string) error {
	cmd := exec.Command("git", "clone", repoPath)
	err := cmd.Run()
	return err
}

func extractGitHubInfo(url string) (string, string, error) {
	// Split the URL by "/"
	parts := strings.Split(url, "/")

	// Check if the URL format is as expected
	if len(parts) == 5 && parts[0] == "https:" && parts[1] == "" && parts[2] == "github.com" {
		username := parts[3]
		repoName := parts[4]
		return username, repoName, nil
	}

	err := errors.New("invalid github url")
	// If the format is not as expected, return an error
	return "", "", err
}
