package env

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
)

const (
	serverPort   = "SERVER_PORT"
	keyStoreFile = "KEY_STORE_FILE"
	passwordFile = "PASSWORD_FILE"

	raftNodeId = "RAFT_NODE_ID"
	raftPort   = "RAFT_PORT"
	raftVolDir = "RAFT_VOL_DIR"
)

var confKeys = []string{
	serverPort,
	keyStoreFile,
	passwordFile,

	raftNodeId,
	raftPort,
	raftVolDir,
}

func GetConf() *Config {
	var v = viper.New()
	v.AutomaticEnv()
	if err := v.BindEnv(confKeys...); err != nil {
		log.Fatal(err)
		return nil
	}

	pass := "111111"
	if v.GetString(passwordFile) != "" {
		passReader, err := ioutil.ReadFile(v.GetString(passwordFile))
		if err != nil {
			fmt.Println("Read passwd file fail", err)
		} else {
			pass = string(passReader)
		}
	}
	privateKey, err := GetPrivateKeyFromKeystore(v.GetString(keyStoreFile), pass)
	if err != nil {
		panic(err)
	}

	conf := Config{
		Server: ConfigServer{
			Port:       v.GetInt(serverPort),
			PrivateKey: privateKey,
		},
		Raft: ConfigRaft{
			NodeId:    v.GetString(raftNodeId),
			Port:      v.GetInt(raftPort),
			VolumeDir: v.GetString(raftVolDir),
		},
	}

	log.Printf("%+v\n", conf)

	return &conf
}

func GetPrivateKeyFromKeystore(keyFile, pass string) (*ecdsa.PrivateKey, error) {
	keyJson, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Println("Read keystore fail", err)
		return nil, err
	}

	keyWrapper, err := keystore.DecryptKey(keyJson, pass)
	if err != nil {
		fmt.Println("Key decrypt error:")
		return nil, err
	}

	fmt.Printf("From address = %s\n", keyWrapper.Address.String())
	return keyWrapper.PrivateKey, nil
}
