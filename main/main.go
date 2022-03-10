package main

import (
	"fmt"
	"github.com/orznewbie/gotest/x/log"
	"time"
)

func parseRFC3339Time(str string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func main() {
	str := "2022-03-04T11:38:43.5018693Z"
	t, err := parseRFC3339Time(str)
	if err != nil {
		fmt.Println(err)
	}
	log.Info(t)
}
