package main

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/require"
	"go.guoyk.net/snowflake"
	"net/http"
	"testing"
)

func BenchmarkServer(b *testing.B) {
	sf := snowflake.New(zeroTime, 0xaa)
	defer sf.Stop()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.Recover())
	routes(e, sf)

	go e.Start(":17001")
	defer e.Shutdown(context.Background())

	for i := 0; i < b.N; i++ {
		hres, err := http.Get("http://127.0.0.1:17001/snowflake/next_ids?size=10&format=str_dec")
		require.NoError(b, err)
		res := NextIDsRes{}
		dec := json.NewDecoder(hres.Body)
		err = dec.Decode(&res)
		require.NoError(b, err)
		hres.Body.Close()
	}
}

func TestServer(t *testing.T) {
	sf := snowflake.New(zeroTime, 0xaa)
	defer sf.Stop()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	routes(e, sf)

	go e.Start(":17001")
	defer e.Shutdown(context.Background())

	hres, err := http.Get("http://127.0.0.1:17001/snowflake/next_ids?size=10&format=str_hex")
	require.NoError(t, err)
	res := NextIDsRes{}
	dec := json.NewDecoder(hres.Body)
	err = dec.Decode(&res)
	require.NoError(t, err)
	hres.Body.Close()
	t.Logf("Resp: %#v", res)
}
