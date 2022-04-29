package json

import (
	"encoding/json"
	"github.com/orznewbie/gotmpl/pkg/log"
	"testing"
)

type User struct {
	Name string
	Age int
}

type Admin struct {
	Id string
	Password string
}

type Wrapper struct {
	P interface{}
}

func TestJson(t *testing.T) {
	w := Wrapper{P: User{
		Name: "zhangshan",
		Age:  10,
	}}
	byt, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}

	var output Wrapper
	if err := json.Unmarshal(byt, &output); err != nil {
		t.Fatal(err)
	}
	log.Info(output.P.(User))
}
