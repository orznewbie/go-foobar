package main

import (
	"encoding/json"
	"fmt"
)

type Raw string

type Wrapper struct {
	Raw Raw
}

func (m Raw) MarshalJSON() ([]byte, error) {
	return []byte(`{"name":"hello world"}`), nil
}

func (m *Raw) UnmarshalJSON(byt []byte) error {
	return nil
}

func main() {
	w := Wrapper{Raw: "fuck"}
	byt, _ := json.Marshal(w)
	fmt.Println(string(byt))
}

//func main() {
//	sr := new(StatusReposne)
//
//	json.Unmarshal([]byte(input), sr)
//	fmt.Printf("%+v\n", sr)
//
//	js, _ := json.Marshal(sr)
//	fmt.Printf("%s\n", js)
//}
//
//type StatusReposne struct {
//	Result []Status `json:"result"`
//}
//
//type Status struct {
//	Id     int
//	Status string
//}
//
//func (x *StatusReposne) MarshalJSON() ([]byte, error) {
//	var buffer struct {
//		Result map[string]string `json:"result"`
//	}
//	buffer.Result = make(map[string]string)
//	for _, v := range x.Result {
//		buffer.Result[strconv.Itoa(v.Id)] = v.Status
//	}
//	return json.Marshal(&buffer)
//}
//
//func (x *StatusReposne) UnmarshalJSON(b []byte) error {
//	var buffer struct {
//		Result map[string]string `json:"result"`
//	}
//	buffer.Result = make(map[string]string)
//	json.Unmarshal(b, &buffer)
//	for k, v := range buffer.Result {
//		k, _ := strconv.Atoi(k)
//		x.Result = append(x.Result, Status{Id: k, Status: v})
//	}
//	return nil
//}
//
//var input = `{
//  "result": {
//    "0": "done",
//    "1": "incomplete",
//    "2": "completed"
//  }
//}`
