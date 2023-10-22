package main

type CredentialValidator interface {
	Match(content string) ([]CloudCredentials, error)
	Validate(CloudCredentials) bool
}

// Function for scanning content , finding credentails and if found then validing
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
