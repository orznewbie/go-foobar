package grpc

import (
	"context"
	testpb "github.com/orznewbie/gotmpl/api/test"
	"github.com/orznewbie/gotmpl/pkg/log"
	rpccode "google.golang.org/genproto/googleapis/rpc/code"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

type DataServiceImpl struct {
	testpb.UnimplementedDataServiceServer
	dataMu     sync.RWMutex
	dataCenter map[string]*testpb.Data

	log log.Logger
}

func (d *DataServiceImpl) GetFile(ctx context.Context, input *testpb.Input) (*testpb.Data, error) {
	d.log.Debug("GetFile request: ", input)
	start := time.Now()

	var ch = make(chan struct{})
	var file = new(testpb.Data)
	go func() {
		d.dataMu.RLock()
		file = d.dataCenter[input.Name]
		// 模拟文件查询耗时
		time.Sleep(time.Second*3)
		d.dataMu.RUnlock()
		ch <- struct{}{}
	}()

	select {
	case <- ctx.Done():
		d.log.Info("从进入到超时经过: ", time.Now().Sub(start))
		return nil, ctx.Err()
	case <- ch:
		if file != nil {
			return file, nil
		}
		return nil, status.Errorf(codes.NotFound, "file not found: %s", input.Name)
	}
}

func (d *DataServiceImpl) Upload(stream testpb.DataService_UploadServer) error {
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&rpcstatus.Status{
				Code:    int32(rpccode.Code_OK),
				Message: "upload successful.",
				Details: nil,
			})
		}
		if err != nil {
			return err
		}
		d.dataMu.Lock()
		d.dataCenter[data.Id] = data
		d.dataMu.Unlock()
	}
}

func (d *DataServiceImpl) Download(input *testpb.Input, stream testpb.DataService_DownloadServer) error {
	d.dataMu.RLock()
	data, ok := d.dataCenter[input.Name]
	d.dataMu.RUnlock()
	if !ok {
		return status.Errorf(codes.NotFound, "data not found: %s", input.Name)
	}
	// 服务端流不需要主动关闭流，在return之后rpc框架会自动关闭，发送一个EOF信号
	for i := 0; i < 3; i++ {
		if err := stream.Send(data); err != nil {
			return err
		}
	}
	return nil
}

func TestDataServiceImpl(t *testing.T) {
	srv := grpc.NewServer()
	impl := &DataServiceImpl{
		dataMu: sync.RWMutex{},
		dataCenter: map[string]*testpb.Data{
			"file0": {
				Id:      "file0",
				Content: "movie data fragment0",
			},
		},
		log: log.Named("data service"),
	}
	testpb.RegisterDataServiceServer(srv, impl)

	lis, err := net.Listen("tcp", "127.0.0.1:223")
	if err != nil {
		t.Fatal(err)
	}
	log.Info("server listening on 127.0.0.1:233..")

	if err := srv.Serve(lis); err != nil {
		t.Fatal(err)
	}
}
