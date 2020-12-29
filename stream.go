package protostream

import (
	"encoding/binary"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"
)

// Stream is a duplex Protocol Buffers streaming agent.
// Messages can be read from a Stream using its Subscribe method.
// Messages can be written to a Stream's destination using the Stream's Push method.
type Stream struct {
	factory       *Factory
	buffer        []byte
	subscriptions map[int][]HandlerFunc
	outBuffer     []byte
	outMut        *sync.RWMutex
	inMut         *sync.Mutex
}

// HandlerFunc is a callback function accepting a Protocol Buffers message.
type HandlerFunc func(data proto.Message)

func (s *Stream) Read(p []byte) (n int, err error) {
	s.outMut.RLock()
	defer s.outMut.RUnlock()

	readIdx := 0

	for readIdx = 0; readIdx < len(p) && readIdx < len(s.outBuffer); readIdx++ {
		p[readIdx] = s.outBuffer[readIdx]
	}

	s.outBuffer = s.outBuffer[readIdx:]

	return readIdx, nil
}

func (s *Stream) Write(p []byte) (n int, err error) {
	s.inMut.Lock()
	defer s.inMut.Unlock()

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

// Push writes a message to outBuffer.
func (s *Stream) Push(data proto.Message) error {
	serial, err := proto.Marshal(data)

	if err != nil {
		return err
	}

	// Find relevant message ID

	typeID, err := s.factory.getTypeID(data)

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

	// Write sendBuf to outBuffer

	s.outMut.Lock()

	s.outBuffer = append(s.outBuffer, sendBuf...)

	s.outMut.Unlock()

	return nil
}

// Subscribe adds h to notification queue for incoming messages of the same type as message.
func (s *Stream) Subscribe(message proto.Message, h HandlerFunc) error {
	messageID, err := s.factory.getTypeID(message)

	if err != nil {
		return err
	}

	subsArray, ok := s.subscriptions[messageID]

	if !ok {
		s.subscriptions[messageID] = []HandlerFunc{}
		subsArray = s.subscriptions[messageID]
	}

	s.subscriptions[messageID] = append(subsArray, h)

	return nil
}
