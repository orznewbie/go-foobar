package dqlx

import (
	"strings"
)

// Escape safely escape a predicate
func Escape(predicate string) string {
	return "<" + escapeSpecialChars(predicate) + ">"
}

func escapeSpecialChars(predicate string) string {
	escapeCharacters := []string{"^", "}", "|", "{", "\\", ",", "<", ">", "\""}

	for _, char := range escapeCharacters {
		predicate = strings.ReplaceAll(predicate, char, "")
	}

	return predicate
}
