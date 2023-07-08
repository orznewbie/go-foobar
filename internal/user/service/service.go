package service

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	userpb "github.com/orznewbie/go-foobar/api/user"
)

type UserServiceImpl struct {
	userpb.UnimplementedUserServiceServer

	mu     *sync.RWMutex
	users  map[int64]*userpb.User
	lastID int64
}

func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{
		mu: new(sync.RWMutex),
		users: map[int64]*userpb.User{
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

func (u *UserServiceImpl) ListUsers(ctx context.Context, in *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	u.mu.RLock()
	users := make([]*userpb.User, 0, len(u.users))
	for _, user := range u.users {
		users = append(users, user)
	}
	u.mu.RUnlock()

	return &userpb.ListUsersResponse{Users: users}, nil
}

func (u *UserServiceImpl) GetUser(ctx context.Context, in *userpb.GetUserRequest) (*userpb.User, error) {
	u.mu.RLock()
	user, ok := u.users[in.Id]
	u.mu.RUnlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user(id=%v)", in.Id)
	}

	return user, nil
}
func (u *UserServiceImpl) CreateUser(ctx context.Context, in *userpb.CreateUserRequest) (*userpb.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	in.User.Id = u.lastID
	u.users[u.lastID] = in.User
	u.lastID++

	return in.User, nil
}

func (u *UserServiceImpl) DeleteUser(ctx context.Context, in *userpb.DeleteUserRequest) (*emptypb.Empty, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	_, ok := u.users[in.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user[id=%d]", in.Id)
	}
	delete(u.users, in.Id)

	return &emptypb.Empty{}, nil
}
