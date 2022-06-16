package custom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

type KV struct {
	Key   string
	Value json.RawMessage
}

type Vertex struct {
	Uid        string
	Attributes []*KV
}

func (v Vertex) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	buf.WriteString(`"uid":"`)
	buf.WriteString(v.Uid)
	buf.WriteByte('"')
	for _, kv := range v.Attributes {
		buf.WriteByte(',')
		buf.WriteString(`"` + kv.Key + `":`)
		buf.Write(kv.Value)
	}
	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (v *Vertex) UnmarshalJSON(byt []byte) error {
	var m = make(map[string]json.RawMessage)
	if err := json.Unmarshal(byt, &m); err != nil {
		return err
	}
	for key, value := range m {
		if key == "uid" {
			v.Uid = string(value)
		} else {
			v.Attributes = append(v.Attributes, &KV{
				Key:   key,
				Value: value,
			})
		}
	}

	return nil
}

func TestCustomMarshal(t *testing.T) {
	money, err := json.Marshal(10000)
	v := Vertex{
		Uid: "0x123",
		Attributes: []*KV{
			{
				Key:   "name",
				Value: []byte(`"张三"`),
			},
			{
				Key:   "money",
				Value: money,
			},
		},
	}
	byt, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(byt))
}

func TestCustomUnmarshal(t *testing.T) {
	JSON := []byte(`{"uid":"0x123", "name":"张三", "money":10000}`)
	var v = new(Vertex)
	if err := json.Unmarshal(JSON, v); err != nil {
		t.Fatal(err)
	}
	for _, kv := range v.Attributes {
		fmt.Println(kv.Key, string(kv.Value))
	}
}
