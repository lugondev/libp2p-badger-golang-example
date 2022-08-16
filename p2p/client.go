package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/multiformats/go-multiaddr"
	"libp2p-badger/db"
	"strconv"
	"time"
)

func StartClient(host string, pingCount int, badger *db.Badger) {
	fmt.Println("Launching p2p")
	client := CreatePeer("/ip4/0.0.0.0/tcp/9001")
	fmt.Printf("Hello World, my hosts ID is %s\n", client.ID().Pretty())
	key := client.ID().Pretty() + "-" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	fmt.Printf("Key request: %s\n", key)
	ma, err := multiaddr.NewMultiaddr(host)
	if err != nil {
		panic(err)
	}
	peerInfo, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	err = client.Connect(ctx, *peerInfo)
	if err != nil {
		panic(err)
	}
	rpcClient := gorpc.NewClient(client, protocolID)
	numCalls := 0
	var durations []time.Duration
	//betweenPingsSleep := time.Second * 1
	betweenPingsSleep := time.Second * 0

	for numCalls < pingCount {
		var reply PingReply
		var args PingArgs

		b := make([]byte, 64)
		_, err := rand.Read(b)
		if err != nil {
			panic(err)
		}
		//fmt.Println("data random:", hex.EncodeToString(b))
		args.Data = b
		args.Key = key

		time.Sleep(betweenPingsSleep)
		startTime := time.Now()
		err = rpcClient.Call(peerInfo.ID, "PingService", "Ping", args, &reply)
		//fmt.Println(string(reply.Data))

		if err != nil {
			panic(err)
		}
		if args.Key != reply.Key {
			panic("Received wrong key!")
		} else {
			err := badger.SetArr(reply.Key, string(reply.Data))
			if err != nil {
				fmt.Println("error SetArr:", reply.Key)
				panic(err)
			}
		}
		endTime := time.Now()
		diff := endTime.Sub(startTime)
		fmt.Printf("%s (%s): seq=%d time=%s\n", peerInfo.ID.String(), peerInfo.Addrs[0].String(), numCalls+1, diff)
		numCalls += 1
		durations = append(durations, diff)
	}

	totalDuration := int64(0)
	for _, dur := range durations {
		totalDuration = totalDuration + dur.Nanoseconds()
	}
	averageDuration := totalDuration / int64(len(durations))
	fmt.Printf("Average duration for ping reply: %s\n", time.Duration(averageDuration))

}
