package user_service

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/genproto/googleapis/longrunning"

	user_v1 "github.com/orznewbie/go-foobar/api/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func NewUserClient() (user_v1.UserServiceClient, *grpc.ClientConn) {
	cc, err := grpc.Dial(UserServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return user_v1.NewUserServiceClient(cc), cc
}

func TestGetUser(t *testing.T) {
	clt, cc := NewUserClient()
	defer cc.Close()

	user, err := clt.GetUser(context.Background(), &user_v1.GetUserRequest{
		Id:      2,
		GetMask: &fieldmaskpb.FieldMask{Paths: []string{"idd", "name"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(user)
}

func TestUpdateUser(t *testing.T) {
	clt, cc := NewUserClient()
	defer cc.Close()

	user, err := clt.UpdateUser(context.Background(), &user_v1.UpdateUserRequest{
		User: &user_v1.User{
			Id:   1,
			Name: "",
			Age:  100,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"xcvassd"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(user)
}

func TestGetOperation(t *testing.T) {
	cc, err := grpc.Dial(UserServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	clt := longrunning.NewOperationsClient(cc)

	op, err := clt.GetOperation(context.Background(), &longrunning.GetOperationRequest{Name: "a"})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(op)
}
