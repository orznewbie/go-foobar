package context

import (
	"context"
	"github.com/orznewbie/gotmpl/pkg/log"
	"testing"
	"time"
)

func TestWithCancel(t *testing.T) {
	sum, err := Cal(context.TODO(), 100)
	if err != nil {
		t.Fatal(err)
	}
	log.Info(sum)
}

func Cal(ctx context.Context, num int) (int, error) {
	var done = make(chan struct{}, 1)
	var sum = 0
	go func() {
		sum = cal(num, done)
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-done:
		return sum, nil
	}
}

func cal(num int, done chan<- struct{}) int {
	var sum = 0
	for i := 1; i <= num; i++ {
		sum += i
	}
	time.Sleep(time.Second * 2)
	done <- struct{}{}
	return sum
}
