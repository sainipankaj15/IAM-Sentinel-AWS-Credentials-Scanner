package main

type CredentialValidator interface {
	Match(content string) ([]CloudCredentials, error)
	Validate(CloudCredentials) bool
}

// Function to scan content and match credentials
// and validate any credentials found
func ParseCredentials(content string, cv CredentialValidator) ([]CloudCredentials, error) {
	res := []CloudCredentials{}
	cc, err := cv.Match(content)
	if err != nil {
		return nil, errorNoMatch
	}

	for _, c := range cc {
		if cv.Validate(c) {
			res = append(res, c)
		}
	}

	return res, nil
}
