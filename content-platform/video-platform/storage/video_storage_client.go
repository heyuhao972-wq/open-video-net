package storage

import videostorage "video-storage"

type StorageClient struct {
	processor *videostorage.Processor
}

type UploadResult struct {
	StorageID    string
	Chunks       []string
	ManifestPath string
	VideoHash    string
	Timestamp    int64
}

func NewStorageClient(basePath string, chunkSize int) (*StorageClient, error) {
	processor, err := videostorage.NewProcessor(basePath, chunkSize)
	if err != nil {
		return nil, err
	}

	return &StorageClient{
		processor: processor,
	}, nil

}

func (c *StorageClient) Upload(path string) (UploadResult, error) {
	result, err := c.processor.StoreVideo(path)
	if err != nil {
		return UploadResult{}, err
	}

	return UploadResult{
		StorageID:    result.VideoID,
		Chunks:       result.ChunkHashes,
		ManifestPath: result.ManifestPath,
		VideoHash:    result.VideoHash,
		Timestamp:    result.Timestamp,
	}, nil
}

func (c *StorageClient) GetChunk(hash string) ([]byte, error) {
	return c.processor.GetChunk(hash)
}

func (c *StorageClient) ComputeManifestHash(path string, authorPublicKey string, videoHash string, timestamp int64) (string, error) {
	return videostorage.ComputeManifestHash(path, authorPublicKey, videoHash, timestamp)
}

func (c *StorageClient) SetManifestProof(path string, authorPublicKey string, signature string, videoHash string, timestamp int64) error {
	return videostorage.SetManifestProof(path, authorPublicKey, signature, videoHash, timestamp)
}
