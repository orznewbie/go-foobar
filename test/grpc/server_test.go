package grpc

import (
	"context"
	"fmt"
	"github.com/orznewbie/gotmpl/test/grpc/pb"
	"google.golang.org/grpc"
	"io"
	"net"
	"testing"
	"time"
)

type CalculateServiceImpl struct {
	pb.UnimplementedCalculateServiceServer
}

func (c CalculateServiceImpl) Sum(ctx context.Context, input *pb.Input) (*pb.Output, error) {
	var result = 0
	for i := 1; i <= int(input.Num); i++ {
		result += i
	}
	return &pb.Output{Result: int64(result)}, nil
}

func (c CalculateServiceImpl) Multi(stream pb.CalculateService_MultiServer) error {
	var result int64 = 1
	for {
		input, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&pb.Output{Result: result})
			time.Sleep(time.Hour)
			return nil
		}

		if err != nil {
			return err
		}

		result *= int64(input.Num)
	}
}

func (c CalculateServiceImpl) Repeat(input *pb.Input, stream pb.CalculateService_RepeatServer) error {
	for i := 1; i <= int(input.Num); i++ {
		if err := stream.Send(&pb.Output{Result: int64(i * 10)}); err != nil {
			fmt.Println("hello world")
			return err
		}
	}
	return nil
}

func TestCalculateService(t *testing.T) {
	srv := grpc.NewServer()
	impl := CalculateServiceImpl{}
	pb.RegisterCalculateServiceServer(srv, impl)

	lis, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		t.Fatal(err)
	}

	if err := srv.Serve(lis); err != nil {
		t.Fatal(err)
	}

}
