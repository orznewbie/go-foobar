package main

import "fmt"

type A struct {
	Name       string
	Attributes map[string]string
}

func main() {
	var m = map[string]string{
		"1": "x",
		"2": "y",
	}

	for a := range m {
		fmt.Println(a)
	}
}
