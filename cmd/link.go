package main

import (
	"fmt"
	"os"
)

var (
	__VERSION__ = "Unknown"
	__BUILD__   = "1970-01-01 00:00:00"
	__COMMIT__  = ""
)

func init() {
	if len(os.Args) == 2 && os.Args[1] == "-v" {
		fmt.Printf("v%s (%s, %s)\r\n", __VERSION__, __BUILD__, __COMMIT__)
		os.Exit(0)
	}
}
