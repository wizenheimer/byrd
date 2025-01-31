// ./src/pkg/utils/email.go
package utils

import (
	"strings"
	"unicode"
)

// GenerateNameFromEmail extracts a readable name from an email address
func generateNameFromEmail(email string) string {
	// Get the part before @ symbol
	parts := strings.Split(email, "@")
	if len(parts) == 0 {
		return email
	}

	localPart := parts[0]

	// Remove numbers and special characters, keeping only letters and separators
	cleanPart := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || r == '.' || r == '_' || r == '-' {
			return r
		}
		return -1
	}, localPart)

	// Handle common cases where email contains dots or underscores
	nameParts := strings.FieldsFunc(cleanPart, func(r rune) bool {
		return r == '.' || r == '_' || r == '-'
	})

	// If no separators found, try to split by capital letters
	if len(nameParts) == 1 {
		nameParts = splitByCase(cleanPart)
	}

	// Properly capitalize each part and filter empty parts
	var validParts []string
	for _, part := range nameParts {
		if len(part) > 0 {
			validPart := strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
			validParts = append(validParts, validPart)
		}
	}

	return strings.Join(validParts, " ")
}

// splitByCase splits a string by uppercase letters
func splitByCase(s string) []string {
	var parts []string
	var current strings.Builder

	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			parts = append(parts, current.String())
			current.Reset()
		}
		current.WriteRune(r)
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// CleanEmailList normalizes and cleans a list of emails.
// Remove list is a list of emails to remove from the list if they exist.
func CleanEmailList(emails []string, removeList []string) []string {
	// Normalize email
	for i, email := range emails {
		emails[i] = NormalizeEmail(email)
	}

	// Normalize remove list
	for i, email := range removeList {
		removeList[i] = NormalizeEmail(email)
	}

	// Remove List Map
	removeMap := make(map[string]bool)
	for _, email := range removeList {
		removeMap[email] = true
	}

	// Existing emails
	var cleanedEmails []string
	seenEmails := make(map[string]bool)

	// Clean emails
	for _, email := range emails {
		if !removeMap[email] && !seenEmails[email] {
			cleanedEmails = append(cleanedEmails, email)
			seenEmails[email] = true
		}
	}

	// Return cleaned emails
	return cleanedEmails
}
