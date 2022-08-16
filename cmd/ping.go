package main

import (
	"flag"
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"libp2p-badger/db"
	"libp2p-badger/p2p"
	"log"
	"os"
)

func main() {

	var mode string
	var dataDir string
	var hostAddress string
	var count int
	flag.StringVar(&mode, "mode", "", "host or client mode")
	flag.StringVar(&dataDir, "data-dir", "", "data dir")
	flag.StringVar(&hostAddress, "host", "", "address of host to connect to")
	flag.IntVar(&count, "count", 300, "number of pings to make")
	flag.Parse()

	if mode == "" {
		log.Fatal("You need to specify '-mode' to be either 'host' or 'client'")
	}

	if mode == "host" {
		p2p.StartHost()
		return
	}
	if mode == "client" {
		if hostAddress == "" {
			log.Fatal("You need to specify '-host' when running as a client")
		}

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

		badgerHandler := db.NewBadger(badgerDB)
		p2p.StartClient(hostAddress, count, badgerHandler)
		return
	}
	log.Fatal("Mode '" + mode + "' not recognized. It has to be either 'host' or 'client'")

}
