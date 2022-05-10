package main

import "fmt"

func main() {
	var arr = make([]int, 0, 5)
	arr = append(arr, 1, 2, 3, 4, 5)
	fmt.Println(arr)
	copy(arr, arr[1:])
	fmt.Println(arr)
}