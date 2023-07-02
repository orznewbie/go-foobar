package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"

	"github.com/dgraph-io/dgo/v210/protos/api"
)

type User struct {
	Name         string       `json:"name"`
	PersonalInfo PersonalInfo `json:"personalInfo"`
	StudyInfo    []StudyInfo  `json:"studyInfo"`
	DType        string       `json:"dgraph.type"`
}

type StudyInfo struct {
	Education string `json:"education"`
	Degree    string `json:"degree"`
}

type PersonalInfo struct {
	Age          int32    `json:"age"`
	Height       float32  `json:"height"`
	Weight       int32    `json:"weight"`
	BornPosition Position `json:"bornPosition"`
}

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func TestStructSchema(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	schema := `
	name: string .
	personalInfo: uid .
	studyInfo: [uid] .
	type User {
		name
		personalInfo
	}
	age: int .
	height: float .
	weight: int .
	bornPosition: uid .
	x: float .
	y: float .
	education: string .
	degree: string .
	`

	if err := dg.Alter(context.Background(), &api.Operation{Schema: schema}); err != nil {
		t.Fatal(err)
	}
}

func TestInsert(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	user := &User{
		Name: "jack",
		PersonalInfo: PersonalInfo{
			Age:    22,
			Height: 173.5,
			Weight: 155,
			BornPosition: Position{
				X: 132.3,
				Y: 23.7,
			},
		},
		StudyInfo: []StudyInfo{
			{
				Education: "fresh",
				Degree:    "HUST",
			},
			{
				Education: "xx",
				Degree:    "yy",
			},
		},
		DType: "User",
	}
	byt, _ := json.Marshal(user)
	_, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{
		SetJson:   byt,
		CommitNow: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	resp, err := dg.NewTxn().Query(context.Background(), `query {
	  q(func:has(dgraph.type))@filter(uid(0xc447)) {
		uid
		name
		personalInfo {
		  age
		  height
		  weight
		  bornPosition {
			x
			y
		  }
		}
		studyInfo {
          education
          degree
		}
	  }
	}`)
	if err != nil {
		t.Fatal(err)
	}

	str := jsoniter.Get(resp.Json, "q", 0).ToString()
	fmt.Println(str)
	var user User
	err = json.Unmarshal([]byte(str), &user)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(user)
}

func TestGetSchema(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	res, err := dg.NewTxn().Query(context.Background(), `schema(pred: [name, age]) {type tokenizer}`)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(res.Json))
}
