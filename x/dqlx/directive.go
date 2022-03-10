package dqlx

import (
	"fmt"
	"strings"
)

type directiveType string

var (
	facets    directiveType = "@facets"    // Done
	recurse   directiveType = "@recurse"   // Done
	cascade   directiveType = "@cascade"   // Done
	normalize directiveType = "@normalize" // Done
	groupBy   directiveType = "@groupby"
)

type facetsDirective struct {
	values []interface{}
}

func (facetsDirective *facetsDirective) toDQL() (query string, args []interface{}, err error) {
	var values []string
	for i, value := range facetsDirective.values {
		switch valueCast := value.(type) {
		case Filter:
			if i == 0 {
				partQuery, partArgs, err := valueCast.toDQL()
				if err != nil {
					return "", nil, err
				}
				return string(facets) + "(" + partQuery + ")", partArgs, nil
			} else {
				return "", nil, fmt.Errorf("facets accepts only one filter as value")
			}
		case *aliasExpr, *asExpr, *Order:
			partDql, partArgs, err := valueCast.(dqlizer).toDQL()
			if err != nil {
				return "", nil, err
			}
			args = append(args, partArgs...)
			values = append(values, partDql)
		case string:
			values = append(values, escapePredicate(valueCast))
		case nil:
			return "", args, nil
		default:
			return "", nil, fmt.Errorf("facets accepts only dqlizer or string as value, given %T", valueCast)
		}
	}

	query = string(facets)
	if len(values) > 0 {
		query += "(" + strings.Join(values, ",") + ")"
	}

	return query, args, nil
}

func (facetsDirective facetsDirective) GetValues() []interface{} {
	return facetsDirective.values
}

// Facets returns the expression for @facets directive
func Facets(values ...interface{}) *facetsDirective {
	return &facetsDirective{
		values: values,
	}
}

type recurseDirective struct {
	depth int
	loop  bool
}

func (recurseDirective *recurseDirective) toDQL() (query string, args []interface{}, err error) {
	var parts []string
	if recurseDirective.depth > 0 {
		parts = append(parts, "depth: "+symbolValuePlaceholder)
		args = append(args, recurseDirective.depth)
	}

	parts = append(parts, "loop: "+symbolValuePlaceholder)
	args = append(args, recurseDirective.loop)

	query = string(recurse) + "(" + strings.Join(parts, ",") + ")"

	return query, args, nil
}

// Recurse returns an expression for @recurse directive
func Recurse(depth int, loop bool) *recurseDirective {
	return &recurseDirective{
		depth: depth,
		loop:  loop,
	}
}

type cascadeDirective struct {
}

func (cascadeDirective *cascadeDirective) toDQL() (query string, args []interface{}, err error) {
	return string(cascade), args, nil
}

// Cascade returns an expression for @cascade directive
func Cascade() *cascadeDirective {
	return &cascadeDirective{}
}

type normalizeDirective struct {
}

func (normalizeDirective *normalizeDirective) toDQL() (query string, args []interface{}, err error) {
	return string(normalize), args, nil
}

// Normalize returns an expression for @normalize directive
func Normalize() *normalizeDirective {
	return &normalizeDirective{}
}
