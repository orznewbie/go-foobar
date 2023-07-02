package dgraph

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/orznewbie/go-foobar/pkg/log"
)

func dgoClient() (*dgo.Dgraph, *grpc.ClientConn) {
	cc, err := grpc.Dial("192.168.30.58:29080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	dc := api.NewDgraphClient(cc)
	dg := dgo.NewDgraphClient(dc)

	return dg, cc
}

func TestAddNoconflictSchema(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	op := &api.Operation{}
	op.Schema = `
		email: string @noconflict .
	`
	ctx := context.Background()
	if err := dg.Alter(ctx, op); err != nil {
		t.Fatal(err)
	}
}

func TestQueryRDF(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	q := `
	query {
	  q(func:has(dgraph.type)) {
		uid
		expand(_all_)
	  }
	}`

	resp, err := dg.NewTxn().QueryRDF(context.Background(), q)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(resp.Rdf))
}

func TestQuerySchema(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	q := `schema(type: <dtdl:test:Room;1>) {}`
	resp, err := dg.NewTxn().Query(context.Background(), q)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(resp.Json))
}

func TestQuery(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	//q := `
	//query {
	//  q(func:has(dgraph.type),first:2) {
	//	uid
	//  }
	//}`

	q := `query{
	q(func: uid(0x2715)) {
    uid
    <dtdl:test:Space::is_part_of> @filter(uid(0x2719))@facets {
      uid
  		}
  }
}`

	resp, err := dg.NewReadOnlyTxn().Query(context.Background(), q)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(resp.Json))
}

func TestMutate(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	quad := `
	_:user <name> "lisi" .
	_:user <age> "35" .
	_:user <dgraph.type> "User" .
	`

	tx := dg.NewTxn()
	resp, err := tx.Mutate(context.Background(), &api.Mutation{
		SetNquads: []byte(quad),
		CommitNow: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Info(resp)

	resp, err = tx.Mutate(context.Background(), &api.Mutation{
		SetNquads: []byte(quad),
		CommitNow: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	//tx.Commit(context.Background())
	//tx.Discard(context.Background())

	log.Info(resp)

}

func TestTx(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	txA := dg.NewTxn()
	txB := dg.NewTxn()

	resp, err := txA.Mutate(context.TODO(), &api.Mutation{
		SetNquads: []byte(`<0xc351> <name> "aaa" .`),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("txA: ", resp.Txn)

	resp, err = txB.Mutate(context.TODO(), &api.Mutation{
		SetNquads: []byte(`<0xc351> <name> "bbb" .`),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("txB: ", resp.Txn)

	if err := txA.Commit(context.TODO()); err != nil {
		t.Fatal(err)
	}
	if err := txB.Commit(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func TestAlterSchema(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	schema := `
<dtdl:test:Space::yyyyyyyyyyyyyyyyyyyy>: string @index(hash) .
<dtdl:test:Space::capacity>: int @index(int) .
type <dtdl:test:Space;1> {
	<dtdl:core:Metadata::etag>
}`

	if err := dg.Alter(context.Background(), &api.Operation{Schema: schema}); err != nil {
		t.Fatal(err)
	}
}

// 测试AlterSchema在被超时或取消时会不会只修改了部分Schema
// 结论：不会
func TestAlterSchemaCancel(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	var buf bytes.Buffer
	for i := 1; i <= 1000000; i++ {
		buf.WriteString(fmt.Sprintf("pred%d: int .\n", i))
	}

	// 在超时或主动Cancel之后，Alter失败，没有写入任何Predicate，所以Alter可以看做有“原子性”
	//ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	//defer cancel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := dg.Alter(ctx, &api.Operation{
			Schema: buf.String(),
		}); err != nil {
			log.Error(err)
			return
		}
	}()

	select {
	case <-time.After(time.Microsecond * 500):
		cancel()
	}

	<-time.After(time.Second)
}

func TestDeleteEdge(t *testing.T) {
	dg, cc := dgoClient()
	defer cc.Close()

	var dqlNquads = []byte("<0x2714> * * .\n")

	//var nquad = &api.NQuad{
	//	Subject:   "0x2714",
	//	Predicate: "*",
	//	ObjectId:  "*",
	//	Namespace: 0,
	//}

	if _, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{
		//Del:       []*api.NQuad{nquad},
		DelNquads: dqlNquads,
		CommitNow: true,
	}); err != nil {
		t.Fatal(err)
	}
}
