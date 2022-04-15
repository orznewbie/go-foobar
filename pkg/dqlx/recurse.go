package dqlx

import "strings"

type RecurseDirective struct {
	depth int64
	loop  bool
}

func (recurseDirective *RecurseDirective) Dql() (string, []interface{}, error) {
	var (
		parts []string
		args  []interface{}
	)

	if recurseDirective.depth > 0 {
		parts = append(parts, "depth: "+symbolValuePlaceholder)
		args = append(args, recurseDirective.depth)
	}

	parts = append(parts, "loop: "+symbolValuePlaceholder)
	args = append(args, recurseDirective.loop)

	return "@recurse(" + strings.Join(parts, ",") + ")", args, nil
}

// Recurse returns an expression for @recurse directive
func Recurse(depth int64, loop bool) *RecurseDirective {
	return &RecurseDirective{
		depth: depth,
		loop:  loop,
	}
}
