package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

type User struct {
	Name     string    `json:"name"`
	Location *Pos      `json:"location"`
	School   []*School `json:"school"`
}

type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type School struct {
	SchoolName string `json:"school_name"`
	SchoolSize int    `json:"school_size"`
}

func TestSchema(t *testing.T) {
	clt, cc := dgoClient()
	defer cc.Close()

	resp, err := clt.NewTxn().Query(context.Background(), `query{
		q(func:uid(0x5)) @recurse {
			uid
			expand(_all_)
		}
	}`)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(resp.Json))
	var user []User
	if err := json.Unmarshal([]byte(jsoniter.Get(resp.Json, "q").ToString()), &user); err != nil {
		t.Fatal(err)
	}

	fmt.Println(user[0].Name, user[0].Location.X, user[0].Location.Y, user[0].School)
}
