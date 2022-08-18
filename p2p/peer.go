package p2p

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/libp2p/go-libp2p"
	libp2pCrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	"log"
	"strings"
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

func CreatePeer(listenAddr string, privateKeyEc *ecdsa.PrivateKey) host.Host {
	privateKey, _, err := libp2pCrypto.GenerateECDSAKeyPair(strings.NewReader(hexutil.EncodeBig(privateKeyEc.D)))
	if err != nil {
		panic(err)
	}

	h, err := libp2p.New(libp2p.ListenAddrStrings(listenAddr), libp2p.Identity(privateKey))
	if err != nil {
		panic(err)
	}
	return h
}

var protocolID = protocol.ID("/p2p/rpc/ping")
