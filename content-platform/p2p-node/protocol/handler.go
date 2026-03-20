package protocol

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/network"
)

func EncodeMessage(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}

func DecodeMessage(data []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return msg, err
}

func ReadMessage(s network.Stream) (Message, error) {
	line, err := ReadLine(s)
	if err != nil {
		return Message{}, err
	}
	return DecodeMessage([]byte(line))
}

func WriteMessage(s network.Stream, msg Message) error {
	data, err := EncodeMessage(msg)
	if err != nil {
		return err
	}
	return WriteLine(s, string(data))
}
