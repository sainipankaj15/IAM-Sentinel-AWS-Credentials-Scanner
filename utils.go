package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func dumpIntoFile() error {

	jsonData, err := json.Marshal(AlreadyVerifiedCommit)
	if err != nil {
		return err
	}

	username, reponame, _ := extractGitHubInfo(RepoName)

	fileName := username + reponame + `.json`
	file, err := os.Create(fileName)

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func readFromFile() error {

	// jsonData, err := json.Marshal(AlreadyVerifiedCommit)
	// if err != nil {
	// 	return err
	// }

	username, reponame, _ := extractGitHubInfo(RepoName)

	fileName := username + reponame + `.json`

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a JSON decoder
	decoder := json.NewDecoder(file)

	// Decode the JSON data into the map
	if err := decoder.Decode(&AlreadyVerifiedCommit); err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

func diffCommits(latest []string, previousVerfiedCommit []string) []string {

	slice1 := previousVerfiedCommit
	slice2 := latest

	// Create a map to store the elements from slice1
	elementsMap := make(map[string]bool)

	// Iterate through slice1 and mark elements as seen in the map
	for _, element := range slice1 {
		elementsMap[element] = true
	}

	// Create a slice to store the differences
	differences := []string{}

	// Iterate through slice2 and check for elements not in slice1
	for _, element := range slice2 {
		if !elementsMap[element] {
			differences = append(differences, element)
		} else {
			fmt.Println("This commit is already verified in previous session Commit ID : ", element)
		}
	}

	// The "differences" slice now contains elements that are in slice2 but not in slice1
	fmt.Println("Differences commits in Previous session and new session ", differences)
	return differences
}
