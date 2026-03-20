package chunk

import (
	"io"
	"os"
)

func SplitFile(filePath string, chunkSize int) ([]Chunk, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var chunks []Chunk
	index := 0

	for {

		buffer := make([]byte, chunkSize)

		n, err := file.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		data := buffer[:n]

		hash := HashChunk(data)

		chunk := Chunk{
			ID:    hash,
			Index: index,
			Data:  data,
			Hash:  hash,
			Size:  n,
		}

		chunks = append(chunks, chunk)

		index++
	}

	return chunks, nil
}
