package main

import (
	"flag"
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"libp2p-badger/fsm"
	"libp2p-badger/p2p"
	"libp2p-badger/server"
	"libp2p-badger/types"
	"log"
	"os"
)

func main() {

	var mode string
	var dataDir string

	flag.StringVar(&mode, "mode", "", "host or client mode")
	flag.StringVar(&dataDir, "data-dir", "", "data dir")
	flag.Parse()

	if mode == "" {
		log.Fatal("You need to specify '-mode' to be either 'host' or 'client'")
	}

	if mode == "host" {
		p2p.StartHost()
		return
	}
	if mode == "client" {
		badgerOpt := badger.DefaultOptions(dataDir)
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

		conf := types.Config{
			Server: types.ConfigServer{
				Port: 2221,
			},
			Raft: types.ConfigRaft{
				NodeId:    "node1",
				Port:      1111,
				VolumeDir: dataDir,
			},
		}

		raftServer, err := fsm.NewRaft(badgerDB, &conf)
		if err != nil {
			log.Fatal(err)
			return
		}

		srv := server.New(fmt.Sprintf(":%d", conf.Server.Port), badgerDB, raftServer, &conf)
		fmt.Println("Server is starting...")
		if err := srv.Start(); err != nil {
			log.Fatal(err)
		}

		return
	}
	log.Fatal("Mode '" + mode + "' not recognized. It has to be either 'host' or 'client'")

}
