package gocron

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-co-op/gocron"
)

func TestGocron(t *testing.T) {
	s := gocron.NewScheduler(time.Local)

	s.Every(1).Second().LimitRunsTo(1).Do(func() { fmt.Println("hello world") })
	s.StartAsync()

	fmt.Println(len(s.Jobs()))
	for {
	}
}
