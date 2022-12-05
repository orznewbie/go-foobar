package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	etcdcli "go.etcd.io/etcd/client/v3"

	"github.com/orznewbie/go-foobar/pkg/log"
)

func TestPut(t *testing.T) {
	cli, err := etcdcli.New(etcdcli.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := cli.Put(ctx, "name", "huhaolong")
	if err != nil {
		t.Fatal(err)
	}
	log.Info(resp)
}

func TestWatch(t *testing.T) {
	cli, err := etcdcli.New(etcdcli.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), "name", etcdcli.WithRev(16))
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
