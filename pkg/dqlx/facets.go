package dqlx

import (
	"fmt"
	"strings"
)

type FacetsDirective struct {
	values []interface{}
}

func (facetsDirective *FacetsDirective) Dql() (string, []interface{}, error) {
	var (
		values []string
		args   []interface{}
	)
	for i, value := range facetsDirective.values {
		switch cast := value.(type) {
		case Filter:
			if i == 0 {
				partQ, partArgs, err := cast.Dql()
				if err != nil {
					return "", nil, err
				}
				return "@facets(" + partQ + ")", partArgs, nil
			} else {
				return "", nil, fmt.Errorf("facets accepts only one filter as value")
			}
		case *AliasExpr, *AsExpr, *Order:
			partQ, partArgs, err := cast.(Dqlizer).Dql()
			if err != nil {
				return "", nil, err
			}
			values = append(values, partQ)
			args = append(args, partArgs...)
		case string:
			values = append(values, Escape(cast))
		case nil:
			return "", nil, nil
		default:
			return "", nil, fmt.Errorf("facets accepts only dqlizer or string as value, given %T", cast)
		}
	}

	if len(values) > 0 {
		return "@facets(" + strings.Join(values, ",") + ")", args, nil
	}

	return "@facets", nil, nil
}

func (facetsDirective FacetsDirective) GetValues() []interface{} {
	return facetsDirective.values
}

// Facets returns the expression for @facets directive
func Facets(values ...interface{}) *FacetsDirective {
	return &FacetsDirective{
		values: values,
	}
}
