package server

import (
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"github.com/hashicorp/raft"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"libp2p-badger/env"
	"libp2p-badger/p2p"
	"libp2p-badger/server/raft_handler"
	"libp2p-badger/server/store_handler"
	"net/http"
	"strings"
	"time"
)

// srv struct handling server
type srv struct {
	listenAddress string
	echo          *echo.Echo
	raft          *raft.Raft
}

// Start the server
func (s srv) Start() error {
	fmt.Println("start server listenAddr:", s.listenAddress)

	return s.echo.StartServer(&http.Server{
		Addr:         s.listenAddress,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
}

type requestP2P struct {
	Host string `json:"host"`
}

// New return new server
func New(listenAddr string, badger *badger.DB, r *raft.Raft, conf *env.Config) *srv {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	raftHandler := raft_handler.New(r)
	if conf.Server.Mode == "host" {
		e.POST("/raft/join", raftHandler.JoinRaftHandler)
		e.POST("/raft/remove", raftHandler.RemoveRaftHandler)
	}
	e.GET("/raft/stats", raftHandler.StatsRaftHandler)

	e.POST("/p2p/test", func(eCtx echo.Context) error {
		var form = requestP2P{}
		if err := eCtx.Bind(&form); err != nil {
			return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
				"error": fmt.Sprintf("error binding: %s", err.Error()),
			})
		}

		form.Host = strings.TrimSpace(form.Host)
		if form.Host == "" {
			return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
				"error": "host is empty",
			})
		}
		p2p.StartClient(form.Host, r, conf.Server.PrivateKey)
		return eCtx.JSON(http.StatusOK, map[string]interface{}{
			"message": "Here is the p2p status",
			"data":    "",
		})
	})

	// DB server
	storeHandler := store_handler.New(r, badger, conf)
	e.GET("/store/:key", storeHandler.Get)
	if r != nil {
		e.POST("/store", storeHandler.Store)
		e.DELETE("/store/:key", storeHandler.Delete)
	}

	return &srv{
		listenAddress: listenAddr,
		echo:          e,
		raft:          r,
	}
}
