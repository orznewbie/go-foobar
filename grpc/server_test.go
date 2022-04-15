package grpc

import (
	"context"
	"fmt"
	"github.com/orznewbie/gotest/grpc/api/test"
	"google.golang.org/grpc"
	"io"
	"net"
	"testing"
	"time"
)

type CalculateServiceImpl struct {
	test.UnimplementedCalculateServiceServer
}

func (c CalculateServiceImpl) Sum(ctx context.Context, input *test.Input) (*test.Output, error) {
	var result = 0
	for i := 1; i <= int(input.Num); i++ {
		result += i
	}
	return &test.Output{Result: int64(result)}, nil
}

func (c CalculateServiceImpl) Multi(stream test.CalculateService_MultiServer) error {
	var result int64 = 1
	for {
		input, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&test.Output{Result: result})
			time.Sleep(time.Hour)
			return nil
		}

		if err != nil {
			return err
		}

		result *= int64(input.Num)
	}
}

func (c CalculateServiceImpl) Repeat(input *test.Input, stream test.CalculateService_RepeatServer) error {
	for i := 1; i <= int(input.Num); i++ {
		if err := stream.Send(&test.Output{Result: int64(i * 10)}); err != nil {
			fmt.Println("hello world")
			return err
		}
	}
	return nil
}

func TestCalculateService(t *testing.T) {
	srv := grpc.NewServer()
	impl := CalculateServiceImpl{}
	test.RegisterCalculateServiceServer(srv, impl)

	lis, err := net.Listen("tcp", "127.0.0.1:223")
	if err != nil {
		t.Fatal(err)
	}

	if err := srv.Serve(lis); err != nil {
		t.Fatal(err)
	}

}
