package channel

import (
	"github.com/orznewbie/gotmpl/pkg/log"
	"testing"
)

func TestChannel(t *testing.T) {

}

type (
	A struct {
		ch chan int
	}
	B struct {
		ch chan int
	}
)

// Channel比较等同于指针比较，不同的结构体实例的channel比较是不一样的
// context包用到了这一点，通过比较done channel来判断两个context是否是同一个
func TestCmp(t *testing.T) {
	a1, a2 := new(A), new(A)
	a1.ch, a2.ch = make(chan int), make(chan int)
	if a1.ch == a2.ch {
		log.Info("a1.ch == a2.ch")
	} else {
		log.Info("a1.ch != a2.ch")
	}
}
