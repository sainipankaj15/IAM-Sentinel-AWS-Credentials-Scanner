# Credential-Scanner

Credential-Scanner is a GoLang script that scans a Git repository (GitHub/GitLab/Bitbucket) for any embedded, valid AWS IAM keys. It checks the validity of the Access Key ID by invoking a basic API call supported by AWS. This tool helps identify potential security risks where AWS IAM keys may have been hardcoded in application code, which could be visible to multiple organizational members, leading to misuse.

## Prerequisites

To run the AWS IAM Key Validator, you need to have Go installed on your system. You can download and install Go from the official Go website: https://golang.org/dl/

## Usage

Clone the repository locally and run:
        
        go run . <repo_url>

for example:

        go run . https://github.com/abhishek-pingsafe/Devops-Node


The script will attempt to validate the provided Access Key ID by making a basic API call to AWS. The report of the scan will be logged in a `logs` directory under the name `<repo_name>-result.txt` (like `Devops-Node-result.txt`).

## Result Screenshot

The result is stored in a `logs/` directory under the name `<repo_name>-result.txt` (Here `Devops-Node-result.txt`).

<img width="1680" alt="image" src="https://github.com/sreenikethMadgula/credential-scanner/assets/56798332/5fa82c6d-136a-4c72-8743-971d71fd61cd">


## Attempted Enhancements
1. Extensibility

To extend support to validate other cloud credentials, the `CredentialValidator` interface can be implemented for another type (say `gcpValidator`) in a new file (like `aws_validator.go`).

A command line argument can be accepted to initialize a `CredentialValidator` as such:

        go run . <repo_url> <cloud>

for example

        go run . https://github.com/abhishek-pingsafe/Devops-Node aws

The `main()` can have a factory to initialize a `CredentialValidator` based on the cloud passed as command line argument.

2. Faster Execution

Concurrency is implemented in `scanDir()` to concurrently process files in a directory.


## Further enhancements

1. Better logging - add a custom logger
2. Base64 decoding
3. Baseline definition - Offer the capability to define a baseline file, to ignore items
during the next scan that are present in the baseline. The user should be able to
generate the baseline file with the help of this script.
This feature is useful if the user is running the script every week and doesn’t want to
see the same findings again and again.
