# go-libp2p-gorpc ping example

Quick example how to build a ping service with go-libp2p-gorpc

## Usage

Run client:

```shell
# start client
SERVER_PORT=2221 RAFT_NODE_ID=node1 RAFT_PORT=1111 RAFT_VOL_DIR=node1 KEY_STORE_FILE=./test/keystore/UTC--2018-10-11T01-26-58.462416324Z--3a1b3b81ed061581558a81f11d63e03129347437 go run cmd/main.go -mode client
```

And then copy one of the "I'm listening on" addresses. In this example, we use
the `127.0.0.1` one which ends up being:

```
/ip4/127.0.0.1/tcp/9001/p2p/QmRkJy4FA5ztudW7yR9Rf8TRKPb4vfyiaqtP75tTb3B36B
```

Start node (host), run:

```shell
#node1
SERVER_PORT=2223 RAFT_NODE_ID=node3 RAFT_PORT=1113 RAFT_VOL_DIR=node_3_data KEY_STORE_FILE=./test/keystore/UTC--2019-03-11T06-23-44.238608862Z--ecf880e334de65cd32a63b7b7567797ed707583b go run cmd/main.go -mode host
```

```shell
#node2
SERVER_PORT=2222 RAFT_NODE_ID=node2 RAFT_PORT=1112 RAFT_VOL_DIR=node_2_data KEY_STORE_FILE=./test/keystore/UTC--2019-03-11T06-20-19.810771134Z--88525df23a7f1b3b549bcfd997ce8160ac7976a9 go run cmd/main.go -mode host
```

Join node:
```shell
curl --request POST \
  --url http://localhost:2221/raft/join \
  --header 'Content-Type: application/json' \
  --data '{
	"node_id": "node_3", 
	"raft_address": "127.0.0.1:1113"
}'
```

```shell
curl --request POST \
  --url http://localhost:2221/raft/join \
  --header 'Content-Type: application/json' \
  --data '{
	"node_id": "node_2", 
	"raft_address": "127.0.0.1:1112"
}'
```

Test p2p: 
```shell
curl --request POST \
  --url http://localhost:2221/p2p/test \
  --header 'Content-Type: application/json' \
  --data '{
	"host":"/ip4/127.0.0.1/tcp/9000/p2p/QmdxbcKXSRvArsuE1xywQhvYY3FvotWyRbLGVdyPChHX4F"
}'
```

Then get `Key request: {key}`. Example `Key request: QmRkJy4FA5ztudW7yR9Rf8TRKPb4vfyiaqtP75tTb3B36B-1660806288145`

Check key data:
```
curl --request GET \
  --url http://localhost:2221/store/QmRkJy4FA5ztudW7yR9Rf8TRKPb4vfyiaqtP75tTb3B36B-1660806288145
```
