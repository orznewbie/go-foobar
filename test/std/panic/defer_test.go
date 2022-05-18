package panic

import (
	"fmt"
	"testing"
)

func foo() (ans int) {
	// defer在return之前执行
	defer func() {
		fmt.Println(ans) 	// 打印 1
	}()
	ans = 1
	return
}

func TestReturnDefer(t *testing.T) {
	foo()
}

func TestPanicDefer(t *testing.T) {
	// defer在panic抛出异常前执行
	defer func() {
		fmt.Println("program panic unexpectedly")
	}()

	var arr []int
	fmt.Println(arr[0])
}

// 在Go的panic机制中，延迟函数的调用在释放堆栈信息之前