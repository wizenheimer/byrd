package utils

import (
	"errors"

	"github.com/clerk/clerk-sdk-go/v2"
)

func GetClerkUserEmail(clerkUser *clerk.User) (string, error) {
	if clerkUser.PrimaryEmailAddressID == nil {
		return "", errors.New("primary email address not found")
	}

	primaryEmailAddress := *clerkUser.PrimaryEmailAddressID
	for _, email := range clerkUser.EmailAddresses {
		if email.ID == primaryEmailAddress {
			return email.EmailAddress, nil
		}
	}

	return "", errors.New("primary email address not found")
}

func GetClerkUserFullName(clerkUser *clerk.User) string {
	fullName := ""

	if clerkUser.FirstName != nil {
		fullName += *clerkUser.FirstName
	}

	if clerkUser.LastName != nil {
		fullName += " " + *clerkUser.LastName
	}

	if clerkUser.PrimaryEmailAddressID != nil {
		email, _ := GetClerkUserEmail(clerkUser) // nolint: errcheck // error is non fatal and handled in the fallback
		if email != "" {
			fullName = generateNameFromEmail(email)
		}
	}

	if fullName == "" {
		fullName = "User"
	}

	return fullName
}

func GenerateWorkspaceName(clerkUser *clerk.User) string {
	firstName := ""
	if clerkUser.FirstName != nil {
		firstName += *clerkUser.FirstName
	}

	if firstName == "" {
		email, _ := GetClerkUserEmail(clerkUser) // nolint: errcheck // error is non fatal and handled in the fallback
		if email != "" {
			firstName = generateNameFromEmail(email)
		}
	}

	if firstName == "" {
		firstName = "User"
	}

	return firstName + "'s Workspace"
}
