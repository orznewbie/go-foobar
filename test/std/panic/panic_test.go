package panic

import (
	"fmt"
	"sync"
	"testing"
)

func f(x int) {
	fmt.Printf("f(%d)\n", x+0/x) // panics if x == 0
	defer fmt.Printf("defer %d\n", x)
	f(x - 1)
}

func TestCallFuncPanic(t *testing.T) {
	f(3)
}

func TestGoroutinePanic(t *testing.T) {
	// 不能捕获到其他goroutine的panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("get panic error", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("get panic error", err)
			}
			wg.Done()
		}()
		f(3)
		wg.Done()
	}()
	wg.Wait()
}

func g() {
	panic(1)
}

func TestRecover(t *testing.T) {
	defer func() {
		switch p := recover(); p {
		case nil:
		case 1:
			fmt.Println("recover 1")
		}
	}()
	g()
}
