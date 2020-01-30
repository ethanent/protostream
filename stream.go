package protostream

import (
	"encoding/binary"
	"io"
	"reflect"

	"github.com/golang/protobuf/proto"
)

type Stream struct {
	factory       *Factory
	buffer        []byte
	subscriptions map[int][]HandlerFunc
	out           io.Writer
}

type HandlerFunc func(data proto.Message)

func (s *Stream) Out(to io.Writer) {
	s.out = to
}

func (s *Stream) Write(p []byte) (n int, err error) {
	s.buffer = append(s.buffer, p...)

	uvVal, uvLen := binary.Uvarint(s.buffer)

	// Ensure that no error occurred while reading the length AND that the buffer is long enough for the message
	if (uvVal != 0 && uvLen > 0) && len(s.buffer) >= int(uvVal)+uvLen {
		// Read message ID

		msgTypeID, msgTypeIDLen := binary.Uvarint(s.buffer[uvLen:])

		// Check for related message

		relatedMessage, ok := s.factory.messages[int(msgTypeID)]

		if !ok {
			// Unknown message type ID. This will be ignored and will corrupt stream.
			return
		}

		// Get the serialized Protocol Buffer portion of buffer

		serializedData := s.buffer[uvLen+msgTypeIDLen:]

		// Create new message from expected

		dataMsg := reflect.New(reflect.TypeOf(relatedMessage).Elem()).Interface().(proto.Message)

		// Parse into dataMsg

		err := proto.Unmarshal(serializedData, dataMsg)

		if err != nil {
			// Not finished reading buf
			return len(p), nil
		}

		// Clear out parsed message

		s.buffer = s.buffer[uvLen+int(uvVal):]

		// Trigger write for parsing in case something was buffered while message was being parsed

		s.Write([]byte{})

		// Broadcast by calling subscriptions

		relevantSubs, ok := s.subscriptions[int(msgTypeID)]

		if !ok {
			// No one is listening to this
		} else {
			for _, sub := range relevantSubs {
				sub(dataMsg)
			}
		}
	}

	return len(p), nil
}

func (s *Stream) Push(data proto.Message) error {
	serial, err := proto.Marshal(data)

	if err != nil {
		return err
	}

	// Find relevant message ID

	typeID, err := s.factory.GetTypeID(data)

	if err != nil {
		return err
	}

	msgIDBuf := make([]byte, 64)
	msgIDLen := binary.PutUvarint(msgIDBuf, uint64(typeID))

	msgIDBuf = msgIDBuf[:msgIDLen]

	// Append ID to start of serialized data buf

	finalSerial := []byte{}

	finalSerial = append(finalSerial, msgIDBuf...)
	finalSerial = append(finalSerial, serial...)

	// Create final []byte
	sendBuf := make([]byte, 64)

	// Append message length to final []byte first
	uvLen := binary.PutUvarint(sendBuf, uint64(len(finalSerial)))

	sendBuf = sendBuf[:uvLen]

	// Add finalSerial buf to sendBuf

	sendBuf = append(sendBuf, finalSerial...)

	_, err = s.out.Write(sendBuf)

	if err != nil {
		return err
	}

	return nil
}

func (s *Stream) Subscribe(message int, h HandlerFunc) {
	subsArray, ok := s.subscriptions[message]

	if !ok {
		s.subscriptions[message] = []HandlerFunc{}
		subsArray = s.subscriptions[message]
	}

	s.subscriptions[message] = append(subsArray, h)
}
