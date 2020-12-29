# protostream
Protocol Buffer streaming in Go

[![GoDoc](https://godoc.org/github.com/ethanent/protostream?status.svg)](https://godoc.org/github.com/ethanent/protostream)

**Note**: Protostream (not to be confused with another protocol project, Protocore) is under active development and does not yet provide a finalized API.

## Install

```sh
go get github.com/ethanent/protostream
```

## Design Goal

Status: In progress. Design not finalized, this is an early conception.

```go
// For each conn (eg. a net.Conn or a QUIC connection)

wrap := protostream.Wrap(conn)

wrap.Handle(func (*pb.ChatSend) {
    wrap.Send(&pb.Status{
        Code: 1,
    })
})
```
