package p2p

import (
	"crypto/ecdsa"
	"fmt"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/multiformats/go-multiaddr"
	"log"
	"time"
)

func StartHost(privateKey *ecdsa.PrivateKey) {
	log.Println("Launching hostPeer")
	hostPeer := CreatePeer("/ip4/0.0.0.0/tcp/9000", privateKey)

	log.Printf("Hello World, my hosts ID is %s\n", hostPeer.ID().Pretty())
	for _, addr := range hostPeer.Addrs() {
		ipfsAddr, err := multiaddr.NewMultiaddr("/ipfs/" + hostPeer.ID().Pretty())
		if err != nil {
			panic(err)
		}
		peerAddr := addr.Encapsulate(ipfsAddr)
		log.Printf("I'm listening on %s\n", peerAddr)
	}

	rpcHost := gorpc.NewServer(hostPeer, protocolID)

	svc := PingService{}
	err := rpcHost.Register(&svc)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done: host is started")

	for {
		time.Sleep(time.Second * 1)
	}
}
