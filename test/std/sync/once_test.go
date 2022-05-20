package sync

import (
	"fmt"
	"sync"
	"testing"
)

func TestOnce(t *testing.T) {
	var once sync.Once
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(func() {
				fmt.Println("Hello World!")
			})
			wg.Done()
		}()
	}
	wg.Wait()
}
