package protostream

import (
	"os"
	"testing"

	"github.com/ethanent/gochat/pcol"
	"github.com/golang/protobuf/proto"
)

var f1 *Factory

func TestMain(m *testing.M) {
	// Initialize f1 factory

	f1 = NewFactory()

	f1.RegisterMessage(0, &pcol.SendChat{})

	code := m.Run()

	os.Exit(code)
}

func TestTransfer(t *testing.T) {
	// Send many messages to help test durability
	sendCount := 10000

	testMsg := "Hello there! 5"

	s1 := f1.CreateStream()
	s2 := f1.CreateStream()

	s1.Out(s2)
	s2.Out(s1)

	mc := make(chan *pcol.SendChat, sendCount)

	s2.Subscribe(&pcol.SendChat{}, func(msg proto.Message) {
		data := msg.(*pcol.SendChat)

		mc <- data
	})

	go func() {
		for i := 0; i < sendCount; i++ {
			s1.Push(&pcol.SendChat{
				Message: testMsg,
			})
		}
	}()

	for i := 0; i < sendCount; i++ {
		rdat := <-mc

		if rdat.Message != testMsg {
			// This is a problem -- message content has been damaged

			t.Fatal("Unexpected message content")
		}
	}
}
