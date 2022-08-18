package main

import (
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"libp2p-badger/env"
	"libp2p-badger/fsm"
	"libp2p-badger/p2p"
	"libp2p-badger/server"
	"log"
	"os"
)

func main() {
	conf := env.GetConf()
	if conf == nil {
		fmt.Println("Cannot load config")
		return
	}

	badgerOpt := badger.DefaultOptions(conf.Raft.VolumeDir)
	badgerDB, err := badger.Open(badgerOpt)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		if err := badgerDB.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error close badgerDB: %s\n", err.Error())
		}
	}()

	raftServer, err := fsm.NewRaft(badgerDB, conf)
	if err != nil {
		log.Fatal(err)
		return
	}
	if conf.Server.Mode == "host" {
		go p2p.StartHost(conf.Server.PrivateKey)
	}

	srv := server.New(fmt.Sprintf(":%d", conf.Server.Port), badgerDB, raftServer, conf)

	p2p.GetClient(conf.Server.PrivateKey)
	fmt.Println("Server is starting...")
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
