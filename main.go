package main

import (
	"context"
	"errors"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.guoyk.net/env"
	"go.guoyk.net/snowflake"
	"log"
	"net/http"
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

	HealthPath = "/healthz"
)

var (
	optBind      string
	optClusterID uint64
	optWorkerID  uint64

	hostname string

	zeroTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
)

func init() {
	hostname, _ = os.Hostname()
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
	id = id + 1
	return
}

func setup() (err error) {
	if err = env.StringVar(&optBind, "BIND", ":3000"); err != nil {
		return
	}
	if err = env.Uint64Var(&optClusterID, "CLUSTER_ID", 0); err != nil {
		return
	}
	if err = env.Uint64Var(&optWorkerID, "WORKER_ID", 0); err != nil {
		return
	}

	if optClusterID == 0 {
		err = errors.New("CLUSTER_ID not set")
		return
	}

	if optWorkerID == 0 {
		if optWorkerID = extractSequenceID(hostname); optWorkerID == 0 {
			err = errors.New("WORKER_ID not set and hostname contains no sequence id")
			return
		}
	}

	if optClusterID&Uint5Mask != optClusterID {
		err = errors.New("invalid CLUSTER_ID")
		return
	}

	if optWorkerID&Uint5Mask != optWorkerID {
		err = errors.New("invalid WORKER_ID")
		return
	}
	return
}

func main() {
	var err error
	defer exit(&err)

	// setup
	if err = setup(); err != nil {
		return
	}

	// calculate instance id
	instanceId := optClusterID<<5 + optWorkerID

	// log
	log.Printf("starting, bind=%s, cluster_id=%d, worker_id=%d, instance_id=%d", optBind, optClusterID, optWorkerID, instanceId)

	// create snowflake
	sf := snowflake.New(zeroTime, instanceId)
	defer sf.Stop()

	// create echo server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	meter(e)
	route(e, sf)

	// wait start
	chStart := make(chan error, 1)
	go func() {
		chStart <- e.Start(optBind)
	}()

	// wait signal
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err = <-chStart:
		if err != nil {
			log.Printf("failed to start: %s", err.Error())
			return
		}
	case sig := <-chSig:
		log.Printf("caught signal: %s", sig.String())
	}

	// wait few seconds
	time.Sleep(time.Second * 2)

	// shutdown
	if err = e.Shutdown(context.Background()); err != nil {
		return
	}
}

func meter(e *echo.Echo) {
	p := prometheus.NewPrometheus("echo", func(ctx echo.Context) bool {
		return ctx.Path() == HealthPath
	})
	p.Use(e)
}

func route(e *echo.Echo, sf snowflake.Snowflake) {
	e.GET(HealthPath, func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	})
	g := e.Group("/snowflake")
	g.GET("/next_id", func(ctx echo.Context) (err error) {
		// control cache
		ctx.Response().Header().Set("Cache-Control", "no-store")

		// request
		var q NextIDReq
		if err = ctx.Bind(&q); err != nil {
			return
		}

		// response
		var id interface{}
		switch q.Format {
		case "str_oct":
			id = strconv.FormatUint(sf.NewID(), 8)
		case "str_dec":
			id = strconv.FormatUint(sf.NewID(), 10)
		case "str_hex":
			id = strconv.FormatUint(sf.NewID(), 16)
		default:
			id = sf.NewID()
		}
		return ctx.JSON(http.StatusOK, NextIDRes{ID: id})
	})
	g.GET("/next_ids", func(ctx echo.Context) (err error) {
		// control cache
		ctx.Response().Header().Set("Cache-Control", "no-store")

		// request
		var q NextIDsReq
		if err = ctx.Bind(&q); err != nil {
			return
		}

		// response
		var ids interface{}
		switch q.Format {
		case "str_oct":
			out := make([]string, 0, q.Size)
			for i := 0; i < q.Size; i++ {
				out = append(out, strconv.FormatUint(sf.NewID(), 8))
			}
			ids = out
		case "str_dec":
			out := make([]string, 0, q.Size)
			for i := 0; i < q.Size; i++ {
				out = append(out, strconv.FormatUint(sf.NewID(), 10))
			}
			ids = out
		case "str_hex":
			out := make([]string, 0, q.Size)
			for i := 0; i < q.Size; i++ {
				out = append(out, strconv.FormatUint(sf.NewID(), 16))
			}
			ids = out
		default:
			out := make([]uint64, 0, q.Size)
			for i := 0; i < q.Size; i++ {
				out = append(out, sf.NewID())
			}
			ids = out
		}
		err = ctx.JSON(http.StatusOK, NextIDsRes{IDs: ids})
		return
	})
}
