package server

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"libp2p-badger/server/store_handler"
	"net/http"
	"time"
)

// srv struct handling server
type srv struct {
	listenAddress string
	echo          *echo.Echo
}

// Start the server
func (s srv) Start() error {
	return s.echo.StartServer(&http.Server{
		Addr:         s.listenAddress,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
}

// New return new server
func New(listenAddr string, badger *badger.DB) *srv {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	// Store server
	storeHandler := store_handler.New(badger)
	e.POST("/store", storeHandler.Store)
	e.GET("/store/:key", storeHandler.Get)
	e.DELETE("/store/:key", storeHandler.Delete)

	return &srv{
		listenAddress: listenAddr,
		echo:          e,
	}
}
