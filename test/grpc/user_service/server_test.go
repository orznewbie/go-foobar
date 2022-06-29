package user_service

import (
	"context"
	"net"
	"sync"
	"testing"

	testpb "github.com/orznewbie/gotmpl/api/test"
	"github.com/orznewbie/gotmpl/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceImpl struct {
	testpb.UnimplementedUserServiceServer
	mu     *sync.RWMutex
	users  map[uint64]*testpb.User
	lastId uint64
}

func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{
		mu: new(sync.RWMutex),
		users: map[uint64]*testpb.User{
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
		lastId: 4,
	}
}

func (u *UserServiceImpl) GetUser(ctx context.Context, in *testpb.GetUserRequest) (*testpb.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	user, ok := u.users[in.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user(id=%v) not found", in.Id)
	}

	var ret = new(testpb.User)
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
func (u *UserServiceImpl) CreateUser(ctx context.Context, in *testpb.CreateUserRequest) (*testpb.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.users[u.lastId] = in.User
	u.lastId++
	return in.User, nil
}
func (u *UserServiceImpl) UpdateUser(ctx context.Context, in *testpb.UpdateUserRequest) (*testpb.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}

func TestUserServiceServer(t *testing.T) {
	srv := grpc.NewServer()
	impl := NewUserServiceImpl()
	testpb.RegisterUserServiceServer(srv, impl)

	const addr = "127.0.0.1:666"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("server listening on %s..", addr)

	if err := srv.Serve(lis); err != nil {
		t.Fatal(err)
	}
}
