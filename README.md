# protostream
Protocol Buffer streaming in Go

[![GoDoc](https://godoc.org/github.com/ethanent/protostream?status.svg)](https://godoc.org/github.com/ethanent/protostream)

## Install

```sh
go get https://godoc.org/github.com/ethanent/protostream
```

## Usage

```go
// Create a factory

mainFac := protostream.NewFactory()

fac := protostream.NewFactory()

fac.RegisterMessage(0, &pcol.AddUser{})
// Where pcol.Test is a Protocol Buffers message.

// We create two streams for this example:
s1 := fac.CreateStream()
s2 := fac.CreateStream()

// Have streams output to each other.
// In reality, a stream would normally output to something like a net.Conn.
s1.Out(s2)
s2.Out(s1)

// Subscribe to the pcol.AddUser event
s2.Subscribe(&pcol.AddUser{}, func(data proto.Message) {
    // Assert as the correct event type
    parsed := data.(*pcol.AddUser)

    fmt.Println("Got signup for", parsed.Name, "age", parsed.Age)

    // Here, if we wanted to, we could push data from s2 to s1 using s2.Push.
})

// And now we push data using s1. This will output to s2's Write method.
s1.Push(&pcol.AddUser{
    Name: "Ethan",
    Age: 18,
})
```
