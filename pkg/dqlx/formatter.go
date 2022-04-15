package dqlx

import "strings"

var (
	IndentLevel = 1
)

// Minify minifies a dql query
func Minify(query string) string {
	parts := strings.Fields(query)
	return strings.Join(parts, " ")
}

// SetIndentLevel set query string indent level
func SetIndentLevel(level int) {
	IndentLevel = level
}

func indent(times int) string {
	return strings.Repeat("\t", IndentLevel*times)
}
