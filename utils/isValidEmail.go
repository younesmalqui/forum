package utils

import "regexp"

func IsValidEmail(email string) bool {
	// Define the regex pattern
	emailPattern := `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`

	// Compile the regex
	re := regexp.MustCompile(emailPattern)

	// Check if the email matches the regex
	return re.MatchString(email)
}
