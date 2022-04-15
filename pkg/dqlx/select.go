package dqlx

import (
	"fmt"
	"strings"
)

type Fields struct {
	predicates  []interface{}
	indentLevel int
}

func (Fields *Fields) Dql() (string, []interface{}, error) {
	var (
		selectedFields []string
		args           []interface{}
	)

	for _, field := range Fields.predicates {
		switch cast := field.(type) {
		case *CountExpr, *AsExpr, *AliasExpr, *ExpandExpr, *Aggregation:
			computedDql, computedArgs, err := cast.(Dqlizer).Dql()
			if err != nil {
				return "", nil, err
			}
			args = append(args, computedArgs...)
			selectedFields = append(selectedFields, indent(Fields.indentLevel)+computedDql)
		case string:
			selectedFields = append(selectedFields, indent(Fields.indentLevel)+Escape(cast))
		default:
			return "", nil, fmt.Errorf("fields only accept strings or Dqlizer, given %T", cast)
		}
	}

	return strings.Join(selectedFields, "\n"), args, nil
}

// Select adds nodeAttributes to selection set
func Select(predicates ...interface{}) *Fields {
	return &Fields{
		predicates: predicates,
	}
}

type AliasExpr struct {
	aliasName string
	value     interface{}
}

func (aliasExpr *AliasExpr) Dql() (string, []interface{}, error) {
	var (
		value string
		args  []interface{}
	)

	switch valueCast := aliasExpr.value.(type) {
	case *CountExpr, *Aggregation:
		var err error
		value, args, err = valueCast.(Dqlizer).Dql()
		if err != nil {
			return "", nil, err
		}
	case string:
		value = Escape(valueCast)
	default:
		return "", nil, fmt.Errorf("alias only accepts string, count() or aggregation, given %T", valueCast)
	}

	return Escape(aliasExpr.aliasName) + " : " + value, args, nil
}

// Alias allows to alias a field
// Example: dqlx.Query(...).Select(dqlx.Alias("name", "my_name"))
func Alias(aliasName string, value interface{}) *AliasExpr {
	return &AliasExpr{
		aliasName: aliasName,
		value:     value,
	}
}

type AsExpr struct {
	varName string
	value   interface{}
}

func (asExpr *AsExpr) Dql() (string, []interface{}, error) {
	var (
		value string
		args  []interface{}
	)

	switch cast := asExpr.value.(type) {
	case *CountExpr, *Aggregation:
		var err error
		value, args, err = cast.(Dqlizer).Dql()
		if err != nil {
			return "", nil, err
		}
	case string:
		value = Escape(cast)
	default:
		return "", nil, fmt.Errorf("as only accepts string, count() or aggregation, given %T", cast)
	}

	return Escape(asExpr.varName) + " AS " + value, args, nil
}

// As makes a field a variable
// Example: dqlx.Query(...).Select(dqlx.As("C", "a"))
func As(varName string, value interface{}) *AsExpr {
	return &AsExpr{
		varName: varName,
		value:   value,
	}
}

type CountExpr struct {
	predicate string
	others    []interface{}
}

func (countExpr *CountExpr) Dql() (string, []interface{}, error) {
	var builder strings.Builder
	builder.WriteString("count(")
	builder.WriteString(Escape(countExpr.predicate))

	var args []interface{}
	for _, other := range countExpr.others {
		switch cast := other.(type) {
		case *FacetsDirective:
			partQ, partArgs, err := cast.Dql()
			if err != nil {
				return "", nil, err
			}
			args = append(args, partArgs...)
			builder.WriteString(" ")
			builder.WriteString(partQ)
		default:
			return "", nil, fmt.Errorf("count part type not support, given %T", cast)
		}
	}

	builder.WriteString(")")

	return builder.String(), args, nil
}

// String used for getting the count expression
func (countExpr *CountExpr) String() string {
	return "count(" + Escape(countExpr.predicate) + ")"
}

// Count represent the 'count' expression
// Expression: count(predicate)
func Count(predicate string, others ...interface{}) *CountExpr {
	return &CountExpr{
		predicate: predicate,
		others:    others,
	}
}

type ValExpr struct {
	varName string
}

func (valExpr *ValExpr) Dql() (string, []interface{}, error) {
	return "val(" + Escape(valExpr.varName) + ")", nil, nil
}

// String used for getting the val expression
func (valExpr *ValExpr) String() string {
	return "val(" + Escape(valExpr.varName) + ")"
}

// Val returns val expression
// Expression: val(varName)
func Val(varName string) *ValExpr {
	return &ValExpr{
		varName: varName,
	}
}

type ExpandExpr struct {
	value string
}

func (expandExpr *ExpandExpr) Dql() (string, []interface{}, error) {
	if expandExpr.value == "_all_" {
		return "expand(_all_)", nil, nil
	}
	return "expand(" + Escape(expandExpr.value) + ")", nil, nil
}

// Expand returns Expand expression
// Expression: expand(_all_), expand(type)
func Expand(value string) *ExpandExpr {
	return &ExpandExpr{
		value: value,
	}
}
