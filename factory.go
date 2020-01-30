package protostream

import (
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"
)

// Factory is a producer of streams and is aware of message definitions and message IDs.
type Factory struct {
	messages map[int]proto.Message
}

// RegisterMessage saves a message to the Factory for its Streams to use.
// id refers to the id sent to Streams to specify the message type.
func (f *Factory) RegisterMessage(id int, message proto.Message) {
	f.messages[id] = message
}

// CreateStream initializes and returns a Stream using Factory f.
func (f *Factory) CreateStream() *Stream {
	return &Stream{
		factory:       f,
		buffer:        []byte{},
		subscriptions: map[int][]HandlerFunc{},
		out:           nil,
	}
}

func (f *Factory) GetTypeID(msgTest proto.Message) (int, error) {
	var typeID *int = nil

	for id, msg := range f.messages {
		if reflect.TypeOf(msgTest) == reflect.TypeOf(msg) {
			typeID = &id
			break
		}
	}

	if typeID == nil {
		return 0, errors.New("Stream cannot send unregistered message")
	}

	return *typeID, nil
}

// NewFactory initializes and returns a new Factory.
func NewFactory() *Factory {
	return &Factory{
		messages: map[int]proto.Message{},
	}
}
