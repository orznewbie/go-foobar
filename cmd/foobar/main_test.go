package main

import (
	"context"
	"math"
	"net/http"
	"testing"

	grpc_logging "github.com/orznewbie/go-foobar/pkg/grpc-middleware/logging"

	"github.com/orznewbie/go-foobar/pkg/log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	userpb "github.com/orznewbie/go-foobar/api/user"
	"github.com/orznewbie/go-foobar/internal/user/service"
)

const Listen = "0.0.0.0:2230"

func TestServices(t *testing.T) {
	if err := log.SetLogger("",
		"./",
		"foobar.log",
		5,
		20,
		"MB",
		"info",
		1); err != nil {
		t.Fatal(err)
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_logging.PayloadUnaryServerInterceptor(log.Named("server-requests")),
		),
		grpc.MaxRecvMsgSize(math.MaxInt32-1),
		grpc.MaxSendMsgSize(math.MaxInt32-1),
	)
	userImpl := service.NewUserServiceImpl()
	userpb.RegisterUserServiceServer(grpcServer, userImpl)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	err := userpb.RegisterUserServiceHandlerServer(ctx, mux, userImpl)
	if err != nil {
		t.Fatal(err)
	}

	h2s := &http2.Server{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Cookie,Grpc-Timeout,X-Grpc-Web")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			return
		}

		if r.ProtoMajor == 2 {
			grpcServer.ServeHTTP(w, r)
			return
		}

		mux.ServeHTTP(w, r)
	})
	t.Logf("foobar server listen on %s", Listen)
	if err := http.ListenAndServe(Listen, h2c.NewHandler(handler, h2s)); err != nil {
		t.Fatal(err)
	}
}
