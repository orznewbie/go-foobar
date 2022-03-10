package dqlx

var symbolValuePlaceholder = "??"
var symbolNodeTraversal = "->"

// dqlizer implementors are able to define a custom dql statement
type dqlizer interface {
	toDQL() (query string, args []interface{}, err error)
}
