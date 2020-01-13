package main

import (
	"context"
	"errors"
	"flag"
	"go.guoyk.net/nrpc"
	"go.guoyk.net/snowflake"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func exit(err *error) {
	if *err != nil {
		log.Printf("exited with error: %s", (*err).Error())
		os.Exit(1)
	}
	log.Printf("exited")
}

const (
	Uint5Mask = uint64(1<<5) - 1
)

var (
	optBind      string
	optClusterID uint64
	optWorkerID  uint64

	zeroTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
)

func main() {
	var err error
	defer exit(&err)

	flag.StringVar(&optBind, "bind", ":3000", "bind address")
	flag.Uint64Var(&optClusterID, "cluster-id", 0, "cluster id, 5 bits unsigned integer")
	flag.Uint64Var(&optWorkerID, "worker-id", 0, "worker id, 5 bits unsigned integer")
	flag.Parse()

	if optClusterID&Uint5Mask != optClusterID {
		err = errors.New("invalid cluster id")
		return
	}

	if optWorkerID&Uint5Mask != optWorkerID {
		err = errors.New("invalid work id")
		return
	}

	instanceId := optClusterID<<5 + optWorkerID

	sf := snowflake.New(zeroTime, instanceId)
	defer sf.Stop()

	s := nrpc.NewServer(nrpc.ServerOptions{})
	routes(s, sf)

	if err = s.Start(optBind); err != nil {
		return
	}

	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT)
	sig := <-chSig
	log.Printf("caught signal: %s", sig.String())

	s.Shutdown()
}

func routes(s *nrpc.Server, sf snowflake.Snowflake) *nrpc.Server {
	s.HandleFunc("snowflake", "create", func(ctx context.Context, nreq *nrpc.Request, nres *nrpc.Response) (err error) {
		nres.Payload = &CreateResp{ID: int64(sf.NewID())}
		return
	})
	s.HandleFunc("snowflake", "create_s", func(ctx context.Context, nreq *nrpc.Request, nres *nrpc.Response) (err error) {
		nres.Payload = &CreateSResp{ID: strconv.FormatUint(sf.NewID(), 10)}
		return
	})
	s.HandleFunc("snowflake", "batch", func(ctx context.Context, nreq *nrpc.Request, nres *nrpc.Response) (err error) {
		var req BatchReq
		if err = nreq.Unmarshal(&req); err != nil {
			return
		}
		res := BatchResp{IDs: make([]int64, 0, req.Size)}
		for i := 0; i < req.Size; i++ {
			res.IDs = append(res.IDs, int64(sf.NewID()))
		}
		nres.Payload = res
		return
	})
	s.HandleFunc("snowflake", "batch_s", func(ctx context.Context, nreq *nrpc.Request, nres *nrpc.Response) (err error) {
		var req BatchReq
		if err = nreq.Unmarshal(&req); err != nil {
			return
		}
		res := BatchSResp{IDs: make([]string, 0, req.Size)}
		for i := 0; i < req.Size; i++ {
			res.IDs = append(res.IDs, strconv.FormatUint(sf.NewID(), 10))
		}
		nres.Payload = res
		return
	})
	return s
}
