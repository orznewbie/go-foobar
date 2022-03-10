package dqlx

type aggregationType string

var (
	sum aggregationType = "sum" // Done
	min aggregationType = "min" // Done
	max aggregationType = "max" // Done
	avg aggregationType = "avg" // Done
)

type aggregation struct {
	aggregationType aggregationType
	valExpr         *valExpr
}

func (aggregation *aggregation) toDQL() (query string, args []interface{}, err error) {
	valQuery, valArgs, err := aggregation.valExpr.toDQL()
	args = append(args, valArgs)
	query = string(aggregation.aggregationType) + "(" + valQuery + ")"
	return query, args, nil
}

// Sum represent the 'sum' expression
// Expression: sum(valExpr)
func Sum(valExpr *valExpr) *aggregation {
	return &aggregation{
		aggregationType: sum,
		valExpr:         valExpr,
	}
}

// Avg represent the 'avg' expression
// Expression: avg(valExpr)
func Avg(valExpr *valExpr) *aggregation {
	return &aggregation{
		aggregationType: avg,
		valExpr:         valExpr,
	}
}

// Min represent the 'min' expression
// Expression: min(valExpr)
func Min(valExpr *valExpr) *aggregation {
	return &aggregation{
		aggregationType: min,
		valExpr:         valExpr,
	}
}

// Max represent the 'max' expression
// Expression: max(valExpr)
func Max(valExpr *valExpr) *aggregation {
	return &aggregation{
		aggregationType: max,
		valExpr:         valExpr,
	}
}
