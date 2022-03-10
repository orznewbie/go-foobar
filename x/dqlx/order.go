package dqlx

import "fmt"

type orderDirection string

var (
	orderAsc  orderDirection = "orderasc"
	orderDesc orderDirection = "orderdesc"
)

// Order represents order expression
type Order struct {
	orderDirection orderDirection
	value          interface{}
}

func (order Order) toDQL() (query string, args []interface{}, err error) {
	var value string

	switch valueCast := order.value.(type) {
	case *valExpr:
		valQuery, _, err := valueCast.toDQL()
		if err != nil {
			return "", nil, err
		}
		value = valQuery
	case string:
		value = escapePredicate(valueCast)
	default:
		return "", nil, fmt.Errorf("order clause only accept Val() or string predicate, given %v", valueCast)
	}

	return string(order.orderDirection) + ": " + value, args, nil
}

// OrderAsc returns an orderasc expression
func OrderAsc(value interface{}) *Order {
	return &Order{
		orderDirection: orderAsc,
		value:          value,
	}
}

// OrderDesc returns an orderdesc expression
func OrderDesc(value interface{}) *Order {
	return &Order{
		orderDirection: orderDesc,
		value:          value,
	}
}
