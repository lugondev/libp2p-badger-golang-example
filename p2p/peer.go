package p2p

import (
	"context"
	"encoding/hex"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	"log"
)

type PingArgs struct {
	Key  string
	Data []byte
}
type PingReply struct {
	Key  string
	Data []byte
}
type PingService struct{}

func (t *PingService) Ping(_ context.Context, argType PingArgs, replyType *PingReply) error {
	log.Println("Received a Ping call:", hex.EncodeToString(argType.Data))
	replyData := []byte("reply for key " + argType.Key + " with data:" + hex.EncodeToString(argType.Data))
	replyType.Data = replyData
	replyType.Key = argType.Key
	return nil
}

func CreatePeer(listenAddr string) host.Host {
	h, err := libp2p.New(libp2p.ListenAddrStrings(listenAddr))
	if err != nil {
		panic(err)
	}
	return h
}

var protocolID = protocol.ID("/p2p/rpc/ping")
