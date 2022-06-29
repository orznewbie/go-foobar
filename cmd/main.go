package main

import "fmt"

type Filter [][]string

func (f Filter) withDeleted() Filter {
	for i := range f {
		f[i] = append(f[i], "deleted")
	}
	return f
}

type LogOutput int

const (
	LogToStdout LogOutput = iota + 1
	LogToFile
	LogToRemote
)

type Tmp struct {
	A string
	B string
}

func main() {
	var t = Tmp{"a", "b"}
	fmt.Println(LogToStdout)
	var filters = Filter{
		{
			"a",
			"b",
		},
		{
			"c",
			"d",
		},
	}
	fmt.Println(filters.withDeleted())
}
