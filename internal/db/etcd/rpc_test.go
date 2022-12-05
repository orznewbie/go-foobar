package etcd

import (
	"context"
	"fmt"
	"io"
	"testing"

	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/orznewbie/go-foobar/pkg/log"
)

func TestWatchRPC(t *testing.T) {
	cc, err := grpc.Dial("192.168.30.58:2379", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	clt := pb.NewWatchClient(cc)
	stream, err := clt.Watch(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if err := stream.Send(&pb.WatchRequest{
		RequestUnion: &pb.WatchRequest_CreateRequest{CreateRequest: &pb.WatchCreateRequest{
			Key:            []byte("name"),
			RangeEnd:       nil,
			StartRevision:  0,
			ProgressNotify: false,
			PrevKv:         false,
			WatchId:        0,
			Fragment:       false,
		}},
	}); err != nil {
		t.Fatal(err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Info("Watch服务端结束")
			break
		}
		fmt.Println("watch返回", resp)
		if err != nil {
			t.Fatal(err)
		}
		log.Info(resp)
	}
}
