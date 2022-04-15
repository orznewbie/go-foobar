package grpc

import (
	"context"
	"fmt"
	"github.com/orznewbie/gotest/grpc/api/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"testing"
)

func TestSum(t *testing.T) {
	cc, err := grpc.Dial("127.0.0.1:223", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer cc.Close()

	clt := test.NewCalculateServiceClient(cc)

	output, err := clt.Sum(context.TODO(), &test.Input{Num: 15})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(output)
}

func TestMulti(t *testing.T) {
	cc, err := grpc.Dial("127.0.0.1:223", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer cc.Close()

	clt := test.NewCalculateServiceClient(cc)

	stream, err := clt.Multi(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 3; i++ {
		if err := stream.Send(&test.Input{Num: int32(i)}); err != nil {
			t.Fatal(err)
		}
	}

	result, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("客户端", result)
}

func TestRepeat(t *testing.T) {
	cc, err := grpc.Dial("127.0.0.1:223", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer cc.Close()

	clt := test.NewCalculateServiceClient(cc)

	stream, err := clt.Repeat(context.TODO(), &test.Input{Num: int32(10)})
	if err != nil {
		t.Fatal(err)
	}

	for {
		output, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(output)
	}
}
