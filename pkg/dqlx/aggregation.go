package dqlx

type AggregationType string

var (
	sum AggregationType = "sum" // Done
	min AggregationType = "min" // Done
	max AggregationType = "max" // Done
	avg AggregationType = "avg" // Done
)

type Aggregation struct {
	typ     AggregationType
	valExpr *ValExpr
}

func (aggregation *Aggregation) Dql() (string, []interface{}, error) {
	q, args, err := aggregation.valExpr.Dql()
	if err != nil {
		return "", nil, err
	}
	return string(aggregation.typ) + "(" + q + ")", args, nil
}

// Sum represent the 'sum' expression
// Expression: sum(valExpr)
func Sum(valExpr *ValExpr) *Aggregation {
	return &Aggregation{
		typ:     sum,
		valExpr: valExpr,
	}
}

// Avg represent the 'avg' expression
// Expression: avg(valExpr)
func Avg(valExpr *ValExpr) *Aggregation {
	return &Aggregation{
		typ:     avg,
		valExpr: valExpr,
	}
}

// Min represent the 'min' expression
// Expression: min(valExpr)
func Min(valExpr *ValExpr) *Aggregation {
	return &Aggregation{
		typ:     min,
		valExpr: valExpr,
	}
}

// Max represent the 'max' expression
// Expression: max(valExpr)
func Max(valExpr *ValExpr) *Aggregation {
	return &Aggregation{
		typ:     max,
		valExpr: valExpr,
	}
}
