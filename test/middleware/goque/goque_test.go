package goque

import (
	"fmt"
	"testing"

	"github.com/beeker1121/goque"
	"github.com/orznewbie/gotmpl/pkg/log"
)

func TestStack(t *testing.T) {
	s, err := goque.OpenStack("testdata/stack_dir")
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
	q, _ := goque.OpenQueue("testdata/queue_dir")
	q.Enqueue([]byte("xxx"))
	q.EnqueueString("yyy")

	//peek, _ := q.Peek()
	//fmt.Println(peek.ToString())

	q.Dequeue()
	peek, _ := q.Peek()
	fmt.Println(peek.ToString())
	defer q.Close()
}

type (
	Request struct {
		Query string
	}

	Alter struct {
		Schema string
	}
)

func TestObjectQueue(t *testing.T) {
	req := Request{Query: "query{}"}
	alter := Alter{Schema: "type{}"}

	q, _ := goque.OpenQueue("testdata/object_dir")
	//q.EnqueueObject(req)
	//q.EnqueueObject(alter)

	peek, _ := q.Peek()
	peek.ToObject(&req)
	log.Info(req)
	peek, _ = q.Peek()
	peek.ToObject(&alter)
	log.Info(alter)
}
