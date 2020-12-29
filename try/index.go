package main

import (
	"fmt"
	"io"

	"github.com/ethanent/protostream"
	"github.com/ethanent/protostream/try/pcol"
	"github.com/golang/protobuf/proto"
)

func main() {
	fac := protostream.NewFactory()

	fac.RegisterMessage(0, &pcol.Test{})
	fac.RegisterMessage(1, &pcol.TestReceipt{})

	s1 := fac.CreateStream()
	s2 := fac.CreateStream()

	go io.Copy(s2, s1)
	go io.Copy(s1, s2)

	s2.Subscribe(&pcol.Test{Name: "Ethan"}, func(data proto.Message) {
		parsed := data.(*pcol.Test)

		fmt.Println("Got", parsed.Name, "age", parsed.Age)

		s2.Push(&pcol.TestReceipt{
			Ok: true,
		})
	})

	s1.Subscribe(&pcol.TestReceipt{Ok: true}, func(data proto.Message) {
		parsed := data.(*pcol.TestReceipt)

		fmt.Println("Got TestReceipt", parsed.Ok, parsed.Error)
	})

	s1.Push(&pcol.Test{
		Name: "Ethan",
		Age:  18,
	})

	s1.Push(&pcol.Test{
		Name: "Bohn",
		Age:  7,
	})

	s1.Push(&pcol.Test{
		Name: "Test",
		Age:  2576,
	})
}
