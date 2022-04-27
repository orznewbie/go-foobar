package dgraph

import (
	"context"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/orznewbie/gotmpl/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func dgraphClient() (api.DgraphClient, *grpc.ClientConn) {
	cc, err := grpc.Dial("192.168.30.58:29080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	dc := api.NewDgraphClient(cc)

	return dc, cc
}

func TestDgraphLogin(t *testing.T) {
}

func TestDgraphAlter(t *testing.T) {
	dc, cc := dgraphClient()
	defer cc.Close()

	schema := `
	<name>: string @index(hash) .
	<age>: int .
	type <User> {
		<name>
		<age>
	}`

	payload, err := dc.Alter(context.Background(), &api.Operation{
		Schema: schema,
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Info(payload.String())
}

func TestDgraphQuery(t *testing.T) {
	dc, cc := dgraphClient()
	defer cc.Close()

	q := `
	query{
		q(func: has(dgraph.type), first: 2){
			uid
		}
	}`

	resp, err := dc.Query(context.Background(), &api.Request{
		StartTs:    0,
		Query:      q,
		Vars:       nil,
		ReadOnly:   false,
		BestEffort: false,
		CommitNow:  false,
		RespFormat: api.Request_JSON,
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Info(string(resp.Json))
}

// Mutate和Upsert都是调用的Query接口
func TestDgraphMutate(t *testing.T) {
	dc, cc := dgraphClient()
	defer cc.Close()

	quad1 := `
	<0xc351> <name> "zzz" .
	`

	// CommitNow以Request的为准
	// Request.CommitNow=true，则无论Mutation.CommitNow是啥，都会Mutate成功
	// Request.CommitNow=false，则无论Mutation.CommitNow是啥，在Commit之后都会Mutate成功

	// Request.CommitNow=false的事务会被保存在Dgraph，我们可以根据Response.Txn拿到这个事务，在之后进行Commit或Abort
	// Txn.StartTs是事务的开始时间戳，也是事务的唯一标识
	resp, err := dc.Query(context.Background(), &api.Request{
		StartTs:    0,
		BestEffort: false,
		Mutations: []*api.Mutation{
			{
				SetNquads: []byte(quad1),
				CommitNow: true,
			},
			//{
			//	SetNquads:  []byte(quad2),
			//	CommitNow:  false,
			//},
		},
		CommitNow:  false,
		RespFormat: api.Request_JSON,
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Info(resp)
}

func TestDgraphCheckVersion(t *testing.T) {
	dc, cc := dgraphClient()
	defer cc.Close()

	ver, err := dc.CheckVersion(context.Background(), &api.Check{})
	if err != nil {
		t.Fatal(err)
	}

	log.Info(ver)
}

func TestDgraphTx(t *testing.T) {
	dc, cc := dgraphClient()
	defer cc.Close()

	// 未提交的事务可以在这里提交
	// 事务1：81113 xxx
	// 事务2：81115 xxx
	// 事务3：81117 yyy
	txCtx, err := dc.CommitOrAbort(context.Background(), &api.TxnContext{
		StartTs:  81277,
		CommitTs: 0,
		Aborted:  true,
		Keys:     nil,
		Preds:    nil,
		Hash:     "",
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Info(txCtx)
}
