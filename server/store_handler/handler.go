package store_handler

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/hashicorp/raft"
	"libp2p-badger/fsm"
	"libp2p-badger/types"
)

// handler struct handler
type handler struct {
	raft      *raft.Raft
	DbHandler *fsm.BadgerFSM
	DB        *badger.DB
}

func New(raft *raft.Raft, badgerDB *badger.DB, conf *types.Config) *handler {
	badgerHandler := fsm.NewBadger(badgerDB, conf)
	return &handler{
		DbHandler: badgerHandler,
		DB:        badgerDB,
		raft:      raft,
	}
}
