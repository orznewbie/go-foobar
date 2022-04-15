package dqlx

import "fmt"

type Direction string

var (
	orderAsc  Direction = "orderasc"
	orderDesc Direction = "orderdesc"
)

// Order represents order expression
type Order struct {
	direction Direction
	value     interface{}
}

func (order Order) Dql() (string, []interface{}, error) {
	var value string

	switch cast := order.value.(type) {
	case *ValExpr:
		q, _, err := cast.Dql()
		if err != nil {
			return "", nil, err
		}
		value = q
	case string:
		value = Escape(cast)
	default:
		return "", nil, fmt.Errorf("order clause only accept Val() or string predicate, given %v", cast)
	}

	return string(order.direction) + ": " + value, nil, nil
}

// OrderAsc returns an orderasc expression
func OrderAsc(value interface{}) *Order {
	return &Order{
		direction: orderAsc,
		value:     value,
	}
}

// OrderDesc returns an orderdesc expression
func OrderDesc(value interface{}) *Order {
	return &Order{
		direction: orderDesc,
		value:     value,
	}
}
