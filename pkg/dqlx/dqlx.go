package dqlx

var symbolValuePlaceholder = "??"
var symbolNodeTraversal = "->"

// Dqlizer implementors are able to define a custom dql statement
type Dqlizer interface {
	Dql() (string, []interface{}, error)
}
