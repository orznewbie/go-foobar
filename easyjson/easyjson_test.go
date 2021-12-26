package easyjson

import (
	"fmt"
	"testing"
)

func TestMarshalAndUnmarshal(t *testing.T) {
	personBefore := Person{
		Name:    "huhaolong",
		Age:     21,
		School:  "HUST",
		Hobbies: []string{"yummy", "beauty"},
	}
	data, _ := personBefore.MarshalJSON()
	fmt.Println(string(data))

	personAfter := &Person{}
	personAfter.UnmarshalJSON(data)
	fmt.Println(personAfter)
}
