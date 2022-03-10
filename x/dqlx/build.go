package dqlx

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (tree *QueryTree) toDQL() (query string, args []interface{}, err error) {
	var builder strings.Builder

	builder.WriteString(indent(1))
	if tree.varName != "" {
		builder.WriteString(escapePredicate(tree.varName) + " AS ")
	}

	if tree.isVariable {
		builder.WriteString("var")
	} else {
		builder.WriteString("<" + tree.name + ">")
	}

	builder.WriteString("(")
	var queries []string
	// RootFilter
	// If rootFilter is nil, then order and pagination should not appear together
	if tree.rootFilter != nil {
		partQuery, partArgs, err := tree.rootFilter.toDQL()
		if err != nil {
			return "", nil, err
		}
		queries = append(queries, partQuery)
		args = append(args, partArgs...)

		// Order
		for _, order := range tree.order {
			if order != nil {
				partQuery, partArgs, err := order.toDQL()
				if err != nil {
					return "", nil, err
				}
				queries = append(queries, partQuery)
				args = append(args, partArgs...)
			}
		}
		// Pagination
		if tree.wantPagination() {
			partQuery, partArgs, err := tree.pagination.toDQL()
			if err != nil {
				return "", nil, err
			}
			queries = append(queries, partQuery)
			args = append(args, partArgs...)
		}
		builder.WriteString("func: " + strings.Join(queries, ", ") + "")
	}
	builder.WriteString(")")

	// Filter
	if tree.filter != nil {
		partQuery, partArgs, err := tree.filter.toDQL()
		if err != nil {
			return "", nil, err
		}
		builder.WriteString(" @filter(" + partQuery + ")")
		args = append(args, partArgs...)
	}
	// Recurse
	if tree.recurse != nil {
		partQuery, partArgs, err := tree.recurse.toDQL()
		if err != nil {
			return "", nil, err
		}
		builder.WriteString(" ")
		builder.WriteString(partQuery)
		args = append(args, partArgs...)
	}
	// Cascade
	if tree.cascade != nil {
		partQuery, partArgs, err := tree.cascade.toDQL()
		if err != nil {
			return "", nil, err
		}
		builder.WriteString(" ")
		builder.WriteString(partQuery)
		args = append(args, partArgs...)
	}
	// Normalize
	if tree.normalize != nil {
		partQuery, partArgs, err := tree.normalize.toDQL()
		if err != nil {
			return "", nil, err
		}
		builder.WriteString(" ")
		builder.WriteString(partQuery)
		args = append(args, partArgs...)
	}
	builder.WriteString(" {\n")
	// Select
	if tree.hasFields() {
		tree.fields.indentLevel = 2
		partQuery, partArgs, err := tree.fields.toDQL()
		if err != nil {
			return "", nil, err
		}
		builder.WriteString(partQuery)
		args = append(args, partArgs...)
	}

	for _, child := range tree.childQueryTrees["root"] {
		if err := tree.buildRecursively(2, child, &builder, &args); err != nil {
			return "", nil, err
		}
	}

	builder.WriteString("\n")
	builder.WriteString(indent(1) + "}")

	return builder.String(), args, nil
}

func (tree *QueryTree) buildRecursively(depth int, path string, builder *strings.Builder, args *[]interface{}) error {
	node, ok := tree.nodes[path]
	if !ok {
		return nil
	}

	builder.WriteString("\n")
	builder.WriteString(indent(depth))
	if node.varName != "" {
		builder.WriteString(node.varName + " AS ")
	}
	builder.WriteString(node.name)
	// Order
	if node.wantOrder() {
		var orderQueries []string
		for _, order := range node.order {
			if order != nil {
				partQuery, partArgs, err := order.toDQL()
				if err != nil {
					return err
				}
				orderQueries = append(orderQueries, partQuery)
				*args = append(*args, partArgs...)
			}
		}
		builder.WriteString(" (" + strings.Join(orderQueries, ", ") + ")")
	}
	// Pagination
	if node.wantPagination() {
		partQuery, partArgs, err := node.pagination.toDQL()
		if err != nil {
			return err
		}
		*args = append(*args, partArgs...)
		builder.WriteString(" (" + partQuery + ")")
	}
	// Filter
	if node.filter != nil {
		partQuery, partArgs, err := node.filter.toDQL()
		if err != nil {
			return err
		}
		*args = append(*args, partArgs...)
		builder.WriteString(" @filter(" + partQuery + ")")
	}
	// Facets
	var facetQueries []string
	for _, facet := range node.facets {
		if facet != nil {
			partQuery, partArgs, err := facet.toDQL()
			if err != nil {
				return err
			}
			facetQueries = append(facetQueries, partQuery)
			*args = append(*args, partArgs...)
		}
	}
	builder.WriteString(" ")
	builder.WriteString(strings.Join(facetQueries, " "))
	// Select
	hasFields := node.hasFields()
	if hasFields || len(tree.childQueryTrees[path]) > 0 {
		builder.WriteString(" {\n")

		if hasFields {
			node.fields.indentLevel = depth + 1
			fieldQuery, fieldArgs, err := node.fields.toDQL()
			if err != nil {
				return err
			}
			builder.WriteString(fieldQuery)
			*args = append(*args, fieldArgs...)
		}
		for _, child := range tree.childQueryTrees[path] {
			tree.buildRecursively(depth+1, child, builder, args)
		}

		builder.WriteString("\n")
		builder.WriteString(indent(depth) + "}")
	}

	return nil
}

// QueriesToDQL returns the DQL statement for 1 or more queries
// Example: dqlx.QueriesToDQL(query1,query2,query3)
func QueriesToDQL(queries ...*QueryTree) (query string, variables map[string]string, err error) {
	ensureUniqueQueryNames(queries)

	blockNames := make([]string, len(queries))
	for i, query := range queries {
		blockNames[i] = strings.Title(strings.ToLower(query.name))
	}
	queryName := strings.Join(blockNames, "_")

	statements := make([]string, len(queries))
	var args []interface{}
	for i, query := range queries {
		partQuery, partArgs, err := query.toDQL()
		if err != nil {
			return "", nil, err
		}
		statements[i] = partQuery
		args = append(args, partArgs...)
	}

	innerQuery := strings.Join(statements, "\n\n")

	query, rawVariables := replacePlaceholders(innerQuery, args)
	variables, placeholders := toVariables(rawVariables)

	var builder strings.Builder
	builder.WriteString("query " + queryName + "(" + strings.Join(placeholders, ", ") + ") {\n")
	builder.WriteString(query)
	builder.WriteString("\n}")

	return builder.String(), variables, nil
}

func ensureUniqueQueryNames(queries []*QueryTree) {
	queryNames := make(map[string]bool)

	for index, query := range queries {
		if queryNames[query.name] {
			query.Name(query.name + "_" + strconv.Itoa(index))
		}

		queryNames[query.name] = true
	}
}

func replacePlaceholders(queryTree string, args []interface{}) (string, map[int]interface{}) {
	variables := make(map[int]interface{})
	var builder strings.Builder
	i := 0

	for {
		p := strings.Index(queryTree, "??")
		if p == -1 {
			break
		}

		builder.WriteString(queryTree[:p])
		key := "$" + strconv.Itoa(i)
		builder.WriteString(key)
		queryTree = queryTree[p+2:]

		// Assign the variables
		variables[i] = args[i]

		i++
	}

	builder.WriteString(queryTree)

	return builder.String(), variables
}

func toVariables(rawVariables map[int]interface{}) (variables map[string]string, placeholders []string) {
	variables = make(map[string]string)
	placeholders = make([]string, len(rawVariables))

	for i := 0; i < len(rawVariables); i++ {
		variableName := "$" + strconv.Itoa(i)

		variables[variableName] = toVariableValue(rawVariables[i])

		placeholders[i] = variableName + ":" + goTypeToDQLType(rawVariables[i])
	}

	return variables, placeholders
}

func toVariableValue(value interface{}) string {
	switch val := value.(type) {
	case time.Time:
		return val.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func goTypeToDQLType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int8, int32, int64:
		return "int"
	case float32, float64:
		return "float"
	case bool:
		return "bool"
	case time.Time, *time.Time:
		return "datetime"
	}

	return "string"
}
