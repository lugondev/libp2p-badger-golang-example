# go-libp2p-gorpc ping example

Quick example how to build a ping service with go-libp2p-gorpc

This example has two parts, the `host` and the `client`. You can switch between
them with the `-mode` flag that accepts either `host` or `client` as value.

## Usage

Run host:

```shell
# start host
go run cmd/ping.go -mode host -data-dir ./node1
```

And then copy one of the "I'm listening on" addresses. In this example, we use
the `127.0.0.1` one which ends up being:

```
/ip4/127.0.0.1/tcp/9000/p2p/QmeEJ1NFqgziZfAEt8Wz8ygx43r4e7i5RgGYxvHXjZtf9M
```

Run client:

```shell
#start client
go run cmd/ping.go -mode client -data-dir ./node1
```

then:

```shell
curl --request POST \
  --url http://localhost:2221/p2p/test \
  --header 'Content-Type: application/json' \
  --data '{
	"host":"/ip4/127.0.0.1/tcp/9000/p2p/QmeEJ1NFqgziZfAEt8Wz8ygx43r4e7i5RgGYxvHXjZtf9M"
}'
```

Start node, run:

```shell
#node1
SERVER_PORT=2222 RAFT_NODE_ID=node2 RAFT_PORT=1112 RAFT_VOL_DIR=node_2_data go run cmd/main.go
```

```shell
#node2
SERVER_PORT=2223 RAFT_NODE_ID=node3 RAFT_PORT=1113 RAFT_VOL_DIR=node_3_data go run cmd/main.go
```
