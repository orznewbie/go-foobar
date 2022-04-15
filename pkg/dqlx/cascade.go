package dqlx

type CascadeDirective struct {
}

func (cascadeDirective *CascadeDirective) Dql() (string, []interface{}, error) {
	return "@cascade", nil, nil
}

// Cascade returns an expression for @cascade directive
func Cascade() *CascadeDirective {
	return &CascadeDirective{}
}
