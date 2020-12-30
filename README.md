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
mgr := protostream.NewConnMgr()

mgr.Handle(func (*pb.ChatSend) *pb.Status {
    return *pb.Status{
        Code: 1
    }
})

mgr.Handle(func (*pb.Hello) {
    fmt.Println("Got hello.")
})
```
