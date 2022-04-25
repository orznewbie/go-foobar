package goque

import (
	"fmt"
	"github.com/beeker1121/goque"
	"testing"
)

func TestStack(t *testing.T) {
	s, err := goque.OpenStack("stack_dir")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	//s.Pop()
	s.Push([]byte("aaa"))
	s.PushString("bbb")

	fmt.Println(s.Length())
	peek, _ := s.Peek()
	fmt.Println(string(peek.Value))
}

func TestQueue(t *testing.T) {
	q, _ := goque.OpenQueue("queue_dir")
	q.Enqueue([]byte("xxx"))
	q.EnqueueString("yyy")

	//peek, _ := q.Peek()
	//fmt.Println(peek.ToString())

	q.Dequeue()
	peek, _ := q.Peek()
	fmt.Println(peek.ToString())
	defer q.Close()
}
