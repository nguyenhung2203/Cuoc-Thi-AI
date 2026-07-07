package utils

import (
	"strings"
	"unicode/utf8"
)

// TruncateText safely truncates a string to maxLen characters (runes), adding a suffix if truncated.
func TruncateText(text string, maxLen int) string {
	if utf8.RuneCountInString(text) <= maxLen {
		return text
	}
	
	runes := []rune(text)
	return string(runes[:maxLen]) + "\n... [TRUNCATED DUE TO LENGTH LIMIT]"
}

// CleanJSON removes markdown code block formatting (e.g. ```json ... ```) from a string.
// It uses substring extraction to ignore any text before or after the code block.
func CleanJSON(s string) string {
	s = strings.TrimSpace(s)
	
	firstIdx := strings.Index(s, "```")
	if firstIdx != -1 {
		// Find the first newline after the ```
		newlineIdx := strings.Index(s[firstIdx:], "\n")
		var startIdx int
		if newlineIdx != -1 {
			startIdx = firstIdx + newlineIdx + 1
		} else {
			// fallback if it's just ```json{...}```
			startIdx = firstIdx + 7 // roughly len("```json")
			if startIdx >= len(s) {
				startIdx = firstIdx + 3
			}
		}

		lastIdx := strings.LastIndex(s, "```")
		if lastIdx > startIdx {
			s = s[startIdx:lastIdx]
		}
	}
	
	return strings.TrimSpace(s)
}
