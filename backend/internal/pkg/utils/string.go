package utils

import "unicode/utf8"

// TruncateText safely truncates a string to maxLen characters (runes), adding a suffix if truncated.
func TruncateText(text string, maxLen int) string {
	if utf8.RuneCountInString(text) <= maxLen {
		return text
	}
	
	runes := []rune(text)
	return string(runes[:maxLen]) + "\n... [TRUNCATED DUE TO LENGTH LIMIT]"
}
