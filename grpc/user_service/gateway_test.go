package user_service

import (
	"context"
	"net/http"
	"testing"

	"github.com/orznewbie/go-foobar/api/google/longrunning"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	user_v1 "github.com/orznewbie/go-foobar/api/user/v1"
)

const (
	UserServiceHost          = "127.0.0.1:666"
	UserServiceGateway       = ":1666"
	OperationsServiceGateway = ":2666"
)

func TestAllGateway(t *testing.T) {
	go TestUserService(t)
	go TestUserServiceGateway(t)
	go TestOperationsServiceGateway(t)

	var forever <-chan struct{}
	<-forever
}

func TestUserServiceGateway(t *testing.T) {
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := user_v1.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, UserServiceHost, opts)
	if err != nil {
		t.Fatal(err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	http.ListenAndServe(UserServiceGateway, mux)
}

func TestOperationsServiceGateway(t *testing.T) {
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := longrunning.RegisterOperationsHandlerFromEndpoint(context.Background(), mux, UserServiceHost, opts)
	if err != nil {
		t.Fatal(err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	http.ListenAndServe(OperationsServiceGateway, mux)
}
