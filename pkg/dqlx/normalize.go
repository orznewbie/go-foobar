package dqlx

type NormalizeDirective struct {
}

func (normalizeDirective *NormalizeDirective) Dql() (string, []interface{}, error) {
	return "@normalize", nil, nil
}

// Normalize returns an expression for @normalize directive
func Normalize() *NormalizeDirective {
	return &NormalizeDirective{}
}
