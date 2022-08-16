package store_handler

import (
	"github.com/dgraph-io/badger/v2"
	"libp2p-badger/db"
)

// handler struct handler
type handler struct {
	DbHandler *db.Badger
	DB        *badger.DB
}

func New(badgerDB *badger.DB) *handler {
	badgerHandler := db.NewBadger(badgerDB)
	return &handler{
		DbHandler: badgerHandler,
		DB:        badgerDB,
	}
}
