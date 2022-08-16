package main

import (
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"github.com/spf13/viper"
	"libp2p-badger/server"
	"log"
	"os"
)

const (
	serverPort = "SERVER_PORT"
	dataDir    = "DATA_DIR"
)

var confKeys = []string{
	serverPort,
	dataDir,
}

type ConfigServer struct {
	Port int `mapstructure:"port"`
}

type ConfigDB struct {
	NodeId    string `mapstructure:"node_id"`
	VolumeDir string `mapstructure:"volume_dir"`
}

type Config struct {
	Server ConfigServer `mapstructure:"server"`
	DB     ConfigDB     `mapstructure:"db"`
}

// main entry point of application start
// run using CONFIG=config.yaml ./program
func main() {

	var v = viper.New()
	v.AutomaticEnv()
	if err := v.BindEnv(confKeys...); err != nil {
		log.Fatal(err)
		return
	}

	conf := Config{
		Server: ConfigServer{
			Port: v.GetInt(serverPort),
		},
		DB: ConfigDB{
			NodeId:    "",
			VolumeDir: v.GetString(dataDir),
		},
	}

	log.Printf("%+v\n", conf)

	// Preparing badgerDB
	badgerOpt := badger.DefaultOptions(conf.DB.VolumeDir)
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

	srv := server.New(fmt.Sprintf(":%d", conf.Server.Port), badgerDB)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	return
}
