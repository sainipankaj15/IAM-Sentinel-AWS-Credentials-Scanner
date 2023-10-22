package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

const awsKeyPattern = `(?m)(?i)AKIA[0-9A-Z]{16}\s+\S{40}|AWS[0-9A-Z]{38}\s+?\S{40}`


// Struct
type awsValidator struct{}

// Implementatig Match method
func (a awsValidator) Match(content string) ([]CloudCredentials, error) {
	res := []CloudCredentials{}
	regex := regexp.MustCompile(awsKeyPattern)

	matches := regex.FindAllString(string(content), -1)
	for _, match := range matches {
		matchArr := regexp.MustCompile(`[^\S]+`).Split(match, 2)
		res = append(res, CloudCredentials{
			Id:     matchArr[0],
			Secret: matchArr[1],
		})
	}

	return res, nil

}

// Implemenation of Validate
func (a awsValidator) Validate(c CloudCredentials) bool {
	return validateIAMKeys(c.Id, c.Secret)
}

// Implemenation of ValidateIAM keys
func validateIAMKeys(accessKeyID, secretAccessKey string) bool {
	// Create a new AWS session with the IAM keys
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})

	// Create a new iam service client using the session
	svc := iam.New(sess)

	// Basic API call to check the IAM keys' validity
	d, err := svc.ListGroups(&iam.ListGroupsInput{})
	if err != nil {
		// InvalidClientTokenId error occurs for invalid keys.
		return !strings.Contains(err.Error(), "InvalidClientTokenId")
	}

	fmt.Print(d)

	// IAM keys are valid and the role has permission to list groups
	return true
}
