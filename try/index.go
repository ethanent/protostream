package main

import (
	"fmt"

	"github.com/ethanent/protostream"
	"github.com/ethanent/protostream/try/pcol"
	"github.com/golang/protobuf/proto"
)

func main() {
	fac := protostream.NewFactory()

	fac.RegisterMessage(0, &pcol.Test{})

	s1 := fac.CreateStream()
	s2 := fac.CreateStream()

	s1.Out(s2)
	s2.Out(s1)

	s2.Subscribe(0, func(data proto.Message) {
		fmt.Println("Got", data)
	})

	s1.Push(&pcol.Test{
		Name: "Ethan",
		Age:  18,
	})

	s1.Push(&pcol.Test{
		Name: "Bohn",
		Age:  7,
	})
}
