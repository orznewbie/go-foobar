package dqlx

import "strings"

// Cursor represents pagination parameters
type Cursor struct {
	First  int
	Offset int
	After  string
}

func (p Cursor) toDQL() (query string, args []interface{}, err error) {
	var paginationExpressions []string
	if p.First != 0 {
		paginationExpressions = append(paginationExpressions, "first: "+symbolValuePlaceholder)
		args = append(args, p.First)
	}

	if p.Offset != 0 {
		paginationExpressions = append(paginationExpressions, "offset: "+symbolValuePlaceholder)
		args = append(args, p.Offset)
	}

	if p.After != "" {
		paginationExpressions = append(paginationExpressions, "after: "+symbolValuePlaceholder)
		args = append(args, p.After)
	}

	return strings.Join(paginationExpressions, ", "), args, nil
}
