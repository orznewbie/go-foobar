package main

import "fmt"

func main() {
	defer func() {
		fmt.Println("hello world")
	}()

	foo()
}

func foo() {
	var arr []int
	fmt.Println(arr[0])
}