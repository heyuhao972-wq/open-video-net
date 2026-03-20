package chunk

type Chunk struct {
	ID    string
	Index int
	Data  []byte
	Hash  string
	Size  int
}
