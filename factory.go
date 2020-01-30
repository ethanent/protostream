package protostream

import (
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"
)

type Factory struct {
	messages map[int]proto.Message
}

func (f *Factory) RegisterMessage(id int, message proto.Message) {
	f.messages[id] = message
}

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
		}
	}

	if typeID == nil {
		return 0, errors.New("Stream cannot send unregistered message")
	}

	return *typeID, nil
}

func NewFactory() *Factory {
	return &Factory{
		messages: map[int]proto.Message{},
	}
}
