package user_service

import (
	"context"
	"fmt"
	testpb "github.com/orznewbie/gotmpl/api/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"testing"
)

func NewUserClient() (testpb.UserServiceClient, *grpc.ClientConn) {
	cc, err := grpc.Dial("127.0.0.1:666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return testpb.NewUserServiceClient(cc), cc
}

func TestGetUser(t *testing.T) {
	clt, cc := NewUserClient()
	defer cc.Close()

	user, err := clt.GetUser(context.Background(), &testpb.GetUserRequest{
		Id:      2,
		GetMask: &fieldmaskpb.FieldMask{Paths: []string{"id", "name"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(user)
}

func TestUpdateUser(t *testing.T) {
	clt, cc := NewUserClient()
	defer cc.Close()

	user, err := clt.UpdateUser(context.Background(), &testpb.UpdateUserRequest{
		User: &testpb.User{
			Id:   1,
			Name: "",
			Age:  100,
		},
		UpdateMask: nil,
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(user)
}
