# go-libp2p-gorpc ping example

Quick example how to build a ping service with go-libp2p-gorpc

This example has two parts, the `host` and the `client`. You can switch between
them with the `-mode` flag that accepts either `host` or `client` as value.

## Usage

Have two terminal windows open in the `examples/ping` directory. In the first
one, run:

```
$ go run cmd/ping.go -mode host
```

And then copy one of the "I'm listening on" addresses. In this example, we use
the `127.0.0.1` one which ends up being:

```
/ip4/192.168.30.104/tcp/9000/p2p/QmSA1gi65YswUKrTGj3SNsc5Y7xUJS5RW3xcHK8Lncvmtj
```

Now in the second terminal window, run:

```
$ go run cmd/ping.go -mode client -data-dir ./node1 -host /ip4/192.168.30.104/tcp/9000/p2p/QmSA1gi65YswUKrTGj3SNsc5Y7xUJS5RW3xcHK8Lncvmtj
```

Start server, run:
```
$ SERVER_PORT=2221 DATA_DIR=./node1 go run cmd/main.go
```
