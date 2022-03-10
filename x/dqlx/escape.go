package dqlx

import (
	"strings"
)

// escapePredicate safely escape a predicate
// Example: dqlx.EscapePredicate("predicate")
func escapePredicate(predicate string) string {
	return "<" + escapeSpecialChars(predicate) + ">"
}

func escapeSpecialChars(predicate string) string {
	escapeCharacters := []string{"^", "}", "|", "{", "\\", ",", "<", ">", "\""}

	for _, char := range escapeCharacters {
		predicate = strings.ReplaceAll(predicate, char, "")
	}

	return predicate
}
