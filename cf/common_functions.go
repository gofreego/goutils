package cf

import "regexp"

func IsEmail(s string) bool {
	// Regular expression pattern for a basic email validation
	// This is a simple example and may not cover all possible valid email addresses.
	// More complex patterns are available for stricter validation.
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	match, _ := regexp.MatchString(emailRegex, s)
	return match
}

func IsMobile(s string) bool {
	// Regular expression pattern for a basic mobile number validation
	// This is a simple example and may not cover all possible valid formats.
	// You may need to adapt it for your specific use case.
	mobileRegex := `^\+?[\d\-\(\)]{10,}$`

	match, _ := regexp.MatchString(mobileRegex, s)
	return match
}
