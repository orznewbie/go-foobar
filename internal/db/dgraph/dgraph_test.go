package dgraph

import (
	"context"
	"testing"

	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/orznewbie/go-foobar/pkg/log"
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
func testDgraphMutate(t *testing.T) (tx *api.TxnContext) {
	dc, cc := dgraphClient()
	defer cc.Close()

	// CommitNow以Request的为准
	// Request.CommitNow=true，则无论Mutation.CommitNow是啥，都会Mutate成功
	// Request.CommitNow=false，则无论Mutation.CommitNow是啥，在Commit之后都会Mutate成功
	//
	// Request.CommitNow=false的事务会被保存在Dgraph，我们可以根据Response.Txn拿到这个事务，在之后进行Commit或Abort
	// Txn.StartTs是事务的开始时间戳，也是事务的唯一标识
	resp, err := dc.Query(context.Background(), &api.Request{
		StartTs: 0,
		Mutations: []*api.Mutation{
			{
				SetNquads: []byte(`<0xc351> <name> "ttt" .`),
				CommitNow: false,
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

	return resp.Txn
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

	tx1 := testDgraphMutate(t)
	tx2 := testDgraphMutate(t)

	// 未提交的事务可以在这里提交
	// 事务1：85786 bbb
	// 事务2：85788 ttt
	txCtx, err := dc.CommitOrAbort(context.Background(), tx2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("tx2 ok")

	txCtx, err = dc.CommitOrAbort(context.Background(), tx1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("tx1 ok")

	log.Info(txCtx)
}
