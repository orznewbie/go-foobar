package main

import (
	"github.com/orznewbie/gotmpl/pkg/log"
)

func Add(a, b int) int {
	return a + b
}

func Mul(a, b, c float32) float32 {
	return a * b * c
}

type Sort interface {
	Cmp() int
}

type Req struct {
	Sort
	log    log.Logger
	Method interface{}
	Args   []interface{}
}

func main() {
	req := Req{
		Method: Add,
		Args:   []interface{}{1, 2},
		log:    log.Named("req"),
	}
	req.log.Info("Hello World!")
}
