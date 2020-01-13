package main

import (
	"context"
	"errors"
	"go.guoyk.net/nrpc"
	"go.guoyk.net/snowflake"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	envBind      string
	envClusterID uint64
	envWorkerID  uint64
	envBench     bool

	hostname string

	zeroTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
)

func init() {
	hostname, _ = os.Hostname()
}

func envBoolVar(val *bool, key string, defaultVal bool) {
	sVal := strings.ToUpper(os.Getenv(key))
	if strings.HasPrefix(sVal, "T") || strings.HasPrefix(sVal, "Y") || strings.HasPrefix(sVal, "ON") || strings.HasPrefix(sVal, "1") {
		*val = true
	} else if strings.HasPrefix(sVal, "F") || strings.HasPrefix(sVal, "N") || strings.HasPrefix(sVal, "OFF") || strings.HasPrefix(sVal, "0") {
		*val = false
	} else {
		*val = defaultVal
	}
}

func envStringVar(val *string, key string, defaultVal string) {
	*val = os.Getenv(key)
	if len(*val) == 0 {
		*val = defaultVal
	}
}

func envUint64Var(val *uint64, key string, defaultVal uint64) {
	var err error
	sVal := os.Getenv(key)
	if *val, err = strconv.ParseUint(sVal, 10, 64); err != nil {
		*val = defaultVal
	}
}

func extractSequenceID(hostname string) (id uint64) {
	var err error
	splits := strings.Split(hostname, "-")
	if len(splits) == 0 {
		id = 0
		return
	}
	if id, err = strconv.ParseUint(splits[len(splits)-1], 10, 64); err != nil {
		id = 0
		return
	}
	return
}

func main() {
	var err error
	defer exit(&err)

	envStringVar(&envBind, "BIND", ":3000")
	envUint64Var(&envClusterID, "CLUSTER_ID", 0)
	envUint64Var(&envWorkerID, "WORKER_ID", 0)
	envBoolVar(&envBench, "BENCH", false)

	if envClusterID == 0 {
		err = errors.New("CLUSTER_ID not set")
		return
	}

	if envWorkerID == 0 {
		if envWorkerID = extractSequenceID(hostname); envWorkerID == 0 {
			err = errors.New("WORKER_ID not set and hostname contains no sequence id")
			return
		}
	}

	if envClusterID&Uint5Mask != envClusterID {
		err = errors.New("invalid cluster id")
		return
	}

	if envWorkerID&Uint5Mask != envWorkerID {
		err = errors.New("invalid work id")
		return
	}

	log.Printf("BIND: %s", envBind)
	log.Printf("CLUSTER_ID: %d", envClusterID)
	log.Printf("WORKER_ID: %d", envWorkerID)

	instanceId := envClusterID<<5 + envWorkerID

	sf := snowflake.New(zeroTime, instanceId)
	defer sf.Stop()

	s := nrpc.NewServer(nrpc.ServerOptions{})
	routes(s, sf)

	if err = s.Start(envBind); err != nil {
		return
	}

	if envBench {
		started := time.Now()
		for i := 0; i < 1000; i++ {
			nreq := nrpc.NewRequest("snowflake", "batch")
			nreq.Payload = BatchReq{Size: 10}
			res := BatchResp{}
			if _, err = nrpc.InvokeAddr(context.Background(), "127.0.0.1"+envBind, nreq, &res); err != nil {
				return
			}
		}
		log.Printf("bench for 1000 requests: %s", time.Now().Sub(started).String())
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
