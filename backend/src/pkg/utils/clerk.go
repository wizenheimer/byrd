// ./src/pkg/utils/clerk.go
package utils

import (
	"errors"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/go-petname/petname"
)

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.Trim(email, " "))
}

// GetClerkUserEmail gets the primary email address of a Clerk user.
// It returns an error if the primary email address is not found.
// It returns the primary email address of the Clerk user.
// It normalizes the email address prior to returning it.
func GetClerkUserEmail(clerkUser *clerk.User) (string, error) {
	if clerkUser.PrimaryEmailAddressID == nil {
		return "", errors.New("primary email address not found")
	}

	primaryEmailAddress := *clerkUser.PrimaryEmailAddressID
	for _, email := range clerkUser.EmailAddresses {
		if email.ID == primaryEmailAddress {
			return NormalizeEmail(email.EmailAddress), nil
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

func GenerateWorkspaceName() string {
	workspaceName := petname.Generate(2, "-")
	return workspaceName
}
