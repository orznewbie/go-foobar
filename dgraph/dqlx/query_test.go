package dqlx

import (
	"context"
	"errors"
	"fmt"
	"github.com/fenos/dqlx"
	"strings"
	"testing"
	"unicode"
)

var ErrUnknowFormat = errors.New("unknow format")

type queryer struct {
}

func (q *queryer) Query(ctx context.Context) error {
	panic("implement me")
}

func (q *queryer) Search(ctx context.Context) error {
	panic("implement me")
}

func (q *queryer) Walk(ctx context.Context) error {
	panic("implement me")
}

func (q *queryer) Lookup(on string, filters [][]string, orderBy string, countOnly bool) error {
	dql := dqlx.QueryType(on)

	var or dqlx.Or
	for _, exprs := range filters {
		var and dqlx.And
		for _, expr := range exprs {
			filter, err := q.buildFilter(expr)
			if err != nil {
				return err
			}
			and = append(and, filter)
		}
		or = append(or, and)
	}
	dql = dql.Filter(or)

	orders, err := q.buildOrders(orderBy)
	if err != nil {
		return err
	}
	for _, order := range orders {
		dql = dql.Order(order)
	}

	if countOnly {
		dql = dql.Select(dqlx.Alias("count", dqlx.Count("uid")))
	} else {
		dql = dql.Select(dqlx.RawExpression{Val: "expand(_all_)"})
	}

	fmt.Println(dql.ToDQL())

	return nil
}

func (q *queryer) Pipeline(ctx context.Context) error {
	panic("implement me")
}

// eq(a, "1")
func (q *queryer) buildFilter(expr string) (dqlx.DQLizer, error) {
	args := strings.FieldsFunc(expr, func(c rune) bool {
		return unicode.IsSpace(c) || c == '(' || c == ',' || c == ')'
	})

	if len(args) == 0 {
		return nil, ErrUnknowFormat
	}

	// TODO op validate
	op := args[0]
	fmt.Println(op)

	return dqlx.RawExpression{Val: expr}, nil
}

// foo,bar desc
func (q *queryer) buildOrders(orderBy string) ([]dqlx.DQLizer, error) {
	splits := strings.Split(orderBy, ",")

	var asc, desc []string
	for _, split := range splits {
		args := strings.Fields(split)
		if len(args) == 1 {
			asc = append(asc, args[0])
		} else if len(args) == 2 && args[1] == "desc" {
			desc = append(desc, args[0])
		} else {
			return nil, ErrUnknowFormat
		}
	}

	var orders []dqlx.DQLizer
	for _, by := range asc {
		orders = append(orders, dqlx.OrderAsc(by))
	}
	for _, by := range desc {
		orders = append(orders, dqlx.OrderDesc(by))
	}

	return orders, nil
}

func (q *queryer) buildPagination() {

}

func TestQueryer_Lookup(t *testing.T) {
	q := &queryer{}
	on := "person"
	filters := [][]string{{`eq(name, "alice")`, `gt(age, 20)`}, {`anyofterms(name, "bob charlie")`, `le(age, 20)`}}
	orderBy := "name, age desc"
	countOnly := true
	q.Lookup(on, filters, orderBy, countOnly)
}
