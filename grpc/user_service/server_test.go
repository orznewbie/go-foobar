package user_service

import (
	"context"
	"net"
	"sync"
	"testing"

	"google.golang.org/genproto/googleapis/longrunning"

	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	user_v1 "github.com/orznewbie/go-foobar/api/user/v1"
	"github.com/orznewbie/go-foobar/pkg/log"
)

type UserServiceImpl struct {
	user_v1.UnimplementedUserServiceServer
	longrunning.UnimplementedOperationsServer

	mu     *sync.RWMutex
	users  map[uint64]*user_v1.User
	lastID uint64
}

func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{
		mu: new(sync.RWMutex),
		users: map[uint64]*user_v1.User{
			1: {
				Id:   1,
				Name: "张三",
				Age:  10,
			},
			2: {
				Id:   2,
				Name: "李四",
				Age:  20,
			},
			3: {
				Id:   3,
				Name: "王五",
				Age:  30,
			},
		},
		lastID: 4,
	}
}

func (u *UserServiceImpl) GetUser(ctx context.Context, in *user_v1.GetUserRequest) (*user_v1.User, error) {
	u.mu.RLock()
	user, ok := u.users[in.Id]
	u.mu.RUnlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user(id=%v)", in.Id)
	}

	var ret = new(user_v1.User)
	if in.GetMask == nil {
		ret = user
	} else {
		for _, mask := range in.GetMask.Paths {
			switch mask {
			case "id":
				ret.Id = user.Id
			case "name":
				ret.Name = user.Name
			case "age":
				ret.Age = user.Age
			}
		}
	}

	return ret, nil
}
func (u *UserServiceImpl) CreateUser(ctx context.Context, in *user_v1.CreateUserRequest) (*user_v1.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.users[u.lastID] = in.User
	u.lastID++
	return in.User, nil
}
func (u *UserServiceImpl) UpdateUser(ctx context.Context, in *user_v1.UpdateUserRequest) (*user_v1.User, error) {
	return &user_v1.User{}, nil
	//return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}

func (u *UserServiceImpl) DeleteUser(ctx context.Context, in *user_v1.DeleteUserRequest) (*emptypb.Empty, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	_, ok := u.users[in.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user[id=%d]", in.Id)
	}
	delete(u.users, in.Id)

	return new(emptypb.Empty), nil
}

func TestUserService(t *testing.T) {
	srv := grpc.NewServer()
	impl := NewUserServiceImpl()
	user_v1.RegisterUserServiceServer(srv, impl)
	longrunning.RegisterOperationsServer(srv, impl)

	lis, err := net.Listen("tcp"+
		"", UserServiceHost)
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("server listening on %s..", UserServiceHost)

	if err := srv.Serve(lis); err != nil {
		t.Fatal(err)
	}
}
