package main

import (
	"context"
	"fmt"
	"github.com/orznewbie/gotmpl/pkg/log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	sum, err := cal(ctx, 1000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sum)
	time.Sleep(time.Second * 5)
}

func cal(ctx context.Context, up int) (sum int32, err error) {
	var done = make(chan struct{}, 1)
	go func() {
		for i := 0; i < up; i++ {
			sum += int32(i)
		}
		time.Sleep(time.Second * 2)
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-done:
		return
	}
}
