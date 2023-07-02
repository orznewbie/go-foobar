package user_service

import (
	"fmt"
	"testing"

	user_v1 "github.com/orznewbie/go-foobar/api/user/v1"

	"google.golang.org/protobuf/encoding/protojson"
)

func TestProtoJSON(t *testing.T) {
	byt, err := protojson.Marshal(&user_v1.User{
		Id:   100,
		Name: "xxx",
		Age:  20,
		Role: &user_v1.User_Admin{Admin: &user_v1.Admin{
			Id:  "root",
			Pwd: "123456",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	var user user_v1.User
	if err := protojson.Unmarshal(byt, &user); err != nil {
		t.Fatal(err)
	}
	fmt.Println(user.Role)
}
