package protocol

import (
	"bufio"
	"strings"

	"github.com/libp2p/go-libp2p/core/network"
)

func ReadLine(s network.Stream) (string, error) {
	reader := bufio.NewReader(s)
	msg, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(msg), nil
}

func WriteLine(s network.Stream, line string) error {
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	_, err := s.Write([]byte(line))
	return err
}
