package main

import (
	"errors"
	"fmt"
	"io"

	"log"
	"os"
	"sync"
)

var (
	errorInvalidCreds = errors.New("credentials invalid")
	errorNoMatch      = errors.New("no credentials in content")
)

// Scans a directory for valid secrets;
// performs a scan on each file in the directory
func scanDir(dirName string, out *os.File, cv CredentialValidator) error {
	files, err := os.ReadDir(dirName)
	if err != nil {
		return err
	}

	// Scanning each file in depently in goroutine
	wg := sync.WaitGroup{}
	for _, file := range files {
		if file.IsDir() {
			err = scanDir(dirName+"/"+file.Name(), out, cv)
			if err != nil {
				return err
			}
		} else {
			wg.Add(1)
			go scanFileStart(&wg,file.Name(), dirName, out, cv)
		}

	}

	wg.Wait()
	return nil
}


func scanFileStart(wg *sync.WaitGroup, name string , dirName string, out *os.File, cv CredentialValidator){

	// Early return or panic here would still call wg.Done()
	defer wg.Done()

	file, err := os.Open(dirName + "/" + name)
	if err != nil {
		log.Println(err)
		return
	}

	defer file.Close()

	if err := scanFile(cv, file, out); err != nil {
		switch err {
		case errorInvalidCreds, errorNoMatch:
			return
		default:
			log.Println(err)
			return
		}
	}
}

// Scans a file for valid secrets
func scanFile(cv CredentialValidator, in, out *os.File) error {

	fileContent, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	cc, err := ParseCredentials(string(fileContent), cv)
	if err != nil {
		return err
	}

	result := ""

	for _, value := range cc {
		result += fmt.Sprintf("\t\t\tValid secrets found in file %s:\n\t\t\tAccess Key: %s\n\t\t\tSecret Access Key: %s\n\n", in.Name(), value.Id, value.Secret)
	}

	_, err = out.Write([]byte(result))
	if err != nil {
		return err
	}

	return nil
}
