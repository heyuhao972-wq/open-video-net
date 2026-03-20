package service

import (
	"bytes"
	"io"

	"streaming-service/p2p"
)

type StreamService struct {
	platformBase string
	p2pBase      string
	maxParallel  int
	cacheSize    int
	timeoutMs    int
	retry        int
}

func NewStreamService(platformBase string, p2pBase string, maxParallel int, cacheSize int, timeoutMs int, retry int) *StreamService {

	return &StreamService{
		platformBase: platformBase,
		p2pBase:      p2pBase,
		maxParallel:  maxParallel,
		cacheSize:    cacheSize,
		timeoutMs:    timeoutMs,
		retry:        retry,
	}

}

func (s *StreamService) GetVideoStream(videoId string) (io.Reader, error) {

	fetcher := p2p.NewChunkFetcher(s.platformBase, s.p2pBase, s.maxParallel, s.cacheSize, s.timeoutMs, s.retry)

	chunks, err := fetcher.FetchChunks(videoId)

	if err != nil {
		return nil, err
	}

	buffer := bytes.Buffer{}

	for _, chunk := range chunks {

		buffer.Write(chunk)

	}

	return &buffer, nil

}
