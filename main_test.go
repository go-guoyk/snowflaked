package main

import (
	"context"
	"go.guoyk.net/nrpc"
	"go.guoyk.net/snowflake"
	"testing"
)

func BenchmarkServer(b *testing.B) {
	s := newServer(snowflake.New(zeroTime, 0xaa))
	s.Start(":12201")
	defer s.Shutdown()

	for i := 0; i < b.N; i++ {
		nreq := nrpc.NewRequest("snowflake", "create")
		nreq.Payload = CreateReq{Size: 10}
		resp := CreateResp{}
		_, err := nrpc.Invoke(context.Background(), "127.0.0.1:12201", nreq, &resp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestServer(t *testing.T) {
	s := newServer(snowflake.New(zeroTime, 0xaa))
	s.Start(":12202")
	defer s.Shutdown()

	nreq := nrpc.NewRequest("snowflake", "create")
	nreq.Payload = CreateReq{Size: 10}
	res := CreateResp{}
	nres, err := nrpc.Invoke(context.Background(), "127.0.0.1:12202", nreq, &res)
	t.Logf("%#v", nres)
	t.Logf("%#v", res)
	if err != nil {
		t.Fatal(err)
	}
}
