package main

import (
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

	fmt.Println("All branches are", branches)

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

	msg := fmt.Sprintf("\tBranch: %s\n", branch)
	fmt.Println(msg)
	r.report.Write([]byte(msg))

	commits, err := getAllCommits(dirName)
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
	cmd := exec.Command("git", "-C", "testing", "for-each-ref", "--format=%(refname:short)", "refs/heads/")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
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

// Scans a directory for valid secrets;
// performs a scan on each file in the directory
// func scanDir(dirName string, out *os.File, cv CredentialValidator) error {
// 	files, err := os.ReadDir(dirName)
// 	if err != nil {
// 		return err
// 	}

// 	// scan each file concurrently
// 	wg := sync.WaitGroup{}
// 	for _, file := range files {
// 		if file.IsDir() {
// 			err = scanDir(dirName+"/"+file.Name(), out, cv)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			wg.Add(1)
// 			go func(wg *sync.WaitGroup, name string) {

// 				defer wg.Done()
// 				f, err := os.Open(dirName + "/" + name)
// 				if err != nil {
// 					log.Println(err)
// 					return
// 				}

// 				// close the file after execution if
// 				// there is no error in opening file
// 				defer f.Close()

// 				if err := scanFile(cv, f, out); err != nil {
// 					switch err {
// 					case errorInvalidCreds, errorNoMatch:
// 						return
// 					default:
// 						log.Println(err)
// 						return
// 					}
// 				}

// 			}(&wg, file.Name())
// 		}

// 	}

// 	wg.Wait()
// 	return nil
// }

// // Scans a file for valid secrets
// func scanFile(cv CredentialValidator, in, out *os.File) error {
// 	fileContent, err := io.ReadAll(in)
// 	if err != nil {
// 		return err
// 	}

// 	cc, err := ParseCredentials(string(fileContent), cv)
// 	if err != nil {
// 		return err
// 	}

// 	result := ""
// 	for _, c := range cc {
// 		result += fmt.Sprintf("\t\t\tValid secrets found in file %s:\n\t\t\tAccess Key: %s\n\t\t\tSecret Access Key: %s\n\n", in.Name(), c.Id, c.Secret)
// 	}

// 	_, err = out.Write([]byte(result))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// clone a repository locally; return error if any found
func cloneRepository(repoPath string) error {
	cmd := exec.Command("git", "clone", repoPath)
	err := cmd.Run()
	return err
}
