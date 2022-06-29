package proto_json

import (
	"fmt"
	"testing"

	testpb "github.com/orznewbie/gotmpl/api/test"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestProtoJSON(t *testing.T) {
	byt, err := protojson.Marshal(&testpb.User{
		Id:   100,
		Name: "xxx",
		Age:  20,
		Role: &testpb.User_Admin{Admin: &testpb.Admin{
			Id:  "root",
			Pwd: "123456",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	var user testpb.User
	if err := protojson.Unmarshal(byt, &user); err != nil {
		t.Fatal(err)
	}
	fmt.Println(user.Role)
}
