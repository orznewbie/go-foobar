package data_service

import (
	"context"
	"io"
	"strconv"
	"testing"
	"time"

	testpb "github.com/orznewbie/gotmpl/api/test"
	"github.com/orznewbie/gotmpl/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewDataClient() (testpb.DataServiceClient, *grpc.ClientConn) {
	cc, err := grpc.Dial("127.0.0.1:223", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return testpb.NewDataServiceClient(cc), cc
}

func TestGetFile(t *testing.T) {
	clt, cc := NewDataClient()
	defer cc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	time.Sleep(time.Second * 1)
	file, err := clt.GetFile(ctx, &testpb.Input{Name: "file0"})
	if err != nil {
		t.Fatal(err)
	}

	log.Info(file.Id, file.Content)
}

func TestUpload(t *testing.T) {
	clt, cc := NewDataClient()
	defer cc.Close()

	stream, err := clt.Upload(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if err := stream.Send(&testpb.Data{
			Id:      "file" + strconv.Itoa(i),
			Content: "movie data fragment" + strconv.Itoa(i),
		}); err != nil {
			t.Fatal(err)
		}
	}

	// 服务端流需要主动关闭流，发送一个EOF信号
	result, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
	}

	log.Info(result)
}

func TestDownload(t *testing.T) {
	clt, cc := NewDataClient()
	defer cc.Close()

	stream, err := clt.Download(context.TODO(), &testpb.Input{Name: "file0"})
	if err != nil {
		t.Fatal(err)
	}

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		log.Info(data)
	}
}
