package dqlx

import (
	"fmt"
	"strings"
)

type fields struct {
	predicates  []interface{}
	indentLevel int
}

func (fields *fields) toDQL() (query string, args []interface{}, err error) {
	var selectedFields []string

	for _, field := range fields.predicates {
		switch filedCast := field.(type) {
		case *countExpr, *asExpr, *aliasExpr, *expandExpr:
			computedDql, computedArgs, err := filedCast.(dqlizer).toDQL()
			if err != nil {
				return "", nil, err
			}
			args = append(args, computedArgs...)
			selectedFields = append(selectedFields, indent(fields.indentLevel)+computedDql)
		case string:
			selectedFields = append(selectedFields, indent(fields.indentLevel)+escapePredicate(filedCast))
		default:
			return "", nil, fmt.Errorf("fields only accept strings or dqlizer, given %T", filedCast)
		}
	}

	return strings.Join(selectedFields, "\n"), args, nil
}

// Select adds nodeAttributes to selection set
func Select(predicates ...interface{}) *fields {
	return &fields{
		predicates: predicates,
	}
}

type aliasExpr struct {
	aliasName string
	value     interface{}
}

func (aliasExpr *aliasExpr) toDQL() (query string, args []interface{}, err error) {
	var value string

	switch valueCast := aliasExpr.value.(type) {
	case *countExpr, *aggregation:
		value, args, err = valueCast.(dqlizer).toDQL()
		if err != nil {
			return "", nil, err
		}
	case string:
		value = escapePredicate(valueCast)
	default:
		return "", nil, fmt.Errorf("alias only accepts string, count() or aggregation, given %T", valueCast)
	}

	query = escapePredicate(aliasExpr.aliasName) + " : " + value

	return query, args, nil
}

// Alias allows to alias a field
// Example: dqlx.Query(...).Select(dqlx.Alias("name", "my_name"))
func Alias(aliasName string, value interface{}) *aliasExpr {
	return &aliasExpr{
		aliasName: aliasName,
		value:     value,
	}
}

type asExpr struct {
	varName string
	value   interface{}
}

func (asExpr *asExpr) toDQL() (query string, args []interface{}, err error) {
	var value string

	switch valueCast := asExpr.value.(type) {
	case *countExpr, *aggregation:
		value, args, err = valueCast.(dqlizer).toDQL()
		if err != nil {
			return "", nil, err
		}
	case string:
		value = escapePredicate(valueCast)
	default:
		return "", nil, fmt.Errorf("as only accepts string, count() or aggregation, given %T", valueCast)
	}

	query = escapePredicate(asExpr.varName) + " AS " + value

	return query, args, nil
}

// As makes a field a variable
// Example: dqlx.Query(...).Select(dqlx.As("C", "a"))
func As(varName string, value interface{}) *asExpr {
	return &asExpr{
		varName: varName,
		value:   value,
	}
}

type countExpr struct {
	predicate string
	others    []interface{}
}

func (countExpr *countExpr) toDQL() (query string, args []interface{}, err error) {
	var builder strings.Builder
	builder.WriteString("count(")
	builder.WriteString(escapePredicate(countExpr.predicate))

	for _, other := range countExpr.others {
		switch otherCast := other.(type) {
		case *facetsDirective:
			partQuery, partArgs, err := otherCast.toDQL()
			if err != nil {
				return "", nil, err
			}
			args = append(args, partArgs...)
			builder.WriteString(" ")
			builder.WriteString(partQuery)
		default:
			return "", nil, fmt.Errorf("count part type not support, given %T", otherCast)
		}
	}

	builder.WriteString(")")

	return builder.String(), args, nil
}

// String used for getting the count expression
func (countExpr *countExpr) String() string {
	return "count(" + escapePredicate(countExpr.predicate) + ")"
}

// Count represent the 'count' expression
// Expression: count(predicate)
func Count(predicate string, others ...interface{}) *countExpr {
	return &countExpr{
		predicate: predicate,
		others:    others,
	}
}

type valExpr struct {
	varName string
}

func (valExpr *valExpr) toDQL() (query string, args []interface{}, err error) {
	return "val(" + escapePredicate(valExpr.varName) + ")", args, nil
}

// String used for getting the val expression
func (valExpr *valExpr) String() string {
	return "val(" + escapePredicate(valExpr.varName) + ")"
}

// Val returns val expression
// Expression: val(varName)
func Val(varName string) *valExpr {
	return &valExpr{
		varName: varName,
	}
}

type expandExpr struct {
	value string
}

func (expandExpr *expandExpr) toDQL() (query string, args []interface{}, err error) {
	if expandExpr.value == "_all_" {
		query = "expand(_all_)"
	} else {
		query = "expand(" + escapePredicate(expandExpr.value) + ")"
	}
	return query, args, nil
}

// Expand returns Expand expression
// Expression: expand(_all_), expand(type)
func Expand(value string) *expandExpr {
	return &expandExpr{
		value: value,
	}
}
