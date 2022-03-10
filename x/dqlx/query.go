package dqlx

import (
	"strings"
)

type node struct {
	varName    string
	name       string
	filter     Filter
	pagination *Cursor
	order      []*Order
	facets     []*facetsDirective
	fields     *fields
}

func (node *node) hasFields() bool {
	return node.fields != nil && len(node.fields.predicates) > 0
}

func (node *node) wantOrder() bool {
	return len(node.order) > 0
}

func (node *node) wantPagination() bool {
	return node.pagination != nil && (node.pagination.Offset != 0 || node.pagination.First != 0 || node.pagination.After != "")
}

// QueryTree represents a recursive query tree
type QueryTree struct {
	isVariable bool
	varName    string
	name       string
	rootFilter Filter
	filter     Filter
	recurse    *recurseDirective
	pagination *Cursor
	order      []*Order
	cascade    *cascadeDirective
	normalize  *normalizeDirective
	fields     *fields

	nodes           map[string]*node
	childQueryTrees map[string][]string
}

func (tree *QueryTree) wantPagination() bool {
	return tree.pagination != nil && (tree.pagination.Offset != 0 || tree.pagination.First != 0 || tree.pagination.After != "")
}

func (tree *QueryTree) hasFields() bool {
	return tree.fields != nil && len(tree.fields.predicates) > 0
}

// ToDQL returns the current state of the query as DQL string
// Example: dqlx.Query(...).ToDQL()
func (tree *QueryTree) ToDQL() (query string, args map[string]string, err error) {
	return QueriesToDQL(tree)
}

// Query initialises a query tree with the provided root filter
// example: dqlx.Query(dqlx.Eq(..,..))
func Query(rootFilter Filter) *QueryTree {
	queryTree := &QueryTree{
		isVariable:      false,
		name:            "query",
		rootFilter:      rootFilter,
		nodes:           make(map[string]*node),
		childQueryTrees: make(map[string][]string),
	}
	return queryTree
}

// Name sets the name of the edge
// Example: dqlx.Query(...).Name("bladerunner")
// DQL: { bladerunner(func: ...) { ... }
func (tree *QueryTree) Name(name string) *QueryTree {
	tree.name = name
	return tree
}

// Var initialises a variable query tree
// Example: dqlx.Var(dqlx.Eq(..,..))
func Var(rootFilter Filter) *QueryTree {
	queryTree := Query(rootFilter)
	queryTree.isVariable = true
	queryTree.name = "variable"
	return queryTree
}

// VarAs initialises a variable query tree with a varName
// Example: dqlx.Var(dqlx.Eq(..,..))
func VarAs(varName string, rootFilter Filter) *QueryTree {
	queryTree := Var(rootFilter)
	queryTree.varName = varName
	return queryTree
}

// Filter requests filters for this query
// Example: dqlx.Query(...).Filter(dqlx.Eq{...})
func (tree *QueryTree) Filter(filter Filter) *QueryTree {
	tree.filter = filter
	return tree
}

// Recurse adds recurse directive
func (tree *QueryTree) Recurse(depth int, loop bool) *QueryTree {
	tree.recurse = Recurse(depth, loop)
	return tree
}

// Order requests an ordering for the result set
// Example1: dqlx.Query(...).Order(dqlx.OrderAsc("field1"))
// Example2: dqlx.Query(...).Order(dqlx.OrderDesc("field2"))
func (tree *QueryTree) Order(order ...*Order) *QueryTree {
	tree.order = append(tree.order, order...)
	return tree
}

// OrderAsc alias for ordering in ascending order
// Example:    dqlx.Query(...).OrderAsc("field1")
// Equivalent: dqlx.Query(...).Order(dqlx.OrderAsc("field1"))
func (tree *QueryTree) OrderAsc(predicate interface{}) *QueryTree {
	tree.order = append(tree.order, OrderAsc(predicate))
	return tree
}

// OrderDesc alias for ordering in descending order
// Example:    dqlx.Query(...).OrderDesc("field1")
// Equivalent: dqlx.Query(...).Order(dqlx.OrderDesc("field1"))
func (tree *QueryTree) OrderDesc(predicate interface{}) *QueryTree {
	tree.order = append(tree.order, OrderDesc(predicate))
	return tree
}

// Paginate requests paginated results
// Example: dqlx.Query(...).Paginate(dqlx.Cursor{...})
func (tree *QueryTree) Paginate(pagination Cursor) *QueryTree {
	tree.pagination = &pagination
	return tree
}

// Select assigns predicates to the selection set
// Example: dqlx.Query(...).Select("field1", "field2", "field3")
func (tree *QueryTree) Select(predicates ...interface{}) *QueryTree {
	if len(predicates) == 0 {
		return tree
	}
	tree.fields = Select(predicates...)
	return tree
}

// Cascade represents @cascade directive
func (tree *QueryTree) Cascade() *QueryTree {
	tree.cascade = Cascade()
	return tree
}

// Normalize represents @normalize directive
func (tree *QueryTree) Normalize() *QueryTree {
	tree.normalize = Normalize()
	return tree
}

// Edge adds an edge in the query selection
// Example1: dqlx.Query(...).Edge("path")
// Example2: dqlx.Query(...).Edge("parent->child->child")
// Example3: dqlx.Query(...).Edge("parent->child->child", dqlx.Select(""))
func (tree *QueryTree) Edge(fullPath string, queryParts ...dqlizer) *QueryTree {
	node := &node{}

	last := strings.LastIndex(fullPath, symbolNodeTraversal)
	var father, name string
	if last == -1 {
		father = "root"
		name = fullPath
	} else {
		father = fullPath[:last]
		name = fullPath[last+2:]
	}

	node.name = escapePredicate(name)
	for _, part := range queryParts {
		switch cast := part.(type) {
		case Filter:
			node.filter = cast
		case Cursor:
			node.pagination = &cast
		case *Order:
			node.order = append(node.order, cast)
		case *facetsDirective:
			node.facets = append(node.facets, Facets(cast.GetValues()...))
		case *fields:
			node.fields = Select(cast.predicates...)
		}
	}
	tree.nodes[fullPath] = node
	tree.childQueryTrees[father] = append(tree.childQueryTrees[father], fullPath)

	return tree
}

// EdgeAs adds a new aliased edge
// Example: dqlx.Query(...).EdgeA s("C", "path", ...)
func (tree *QueryTree) EdgeAs(varName string, fullPath string, queryParts ...dqlizer) *QueryTree {
	queryTree := tree.Edge(fullPath, queryParts...)
	queryTree.nodes[fullPath].varName = varName
	return tree
}
