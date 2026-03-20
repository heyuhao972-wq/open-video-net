package storage

import videostorage "video-storage"

type StorageClient struct {
	processor *videostorage.Processor
}

type UploadResult struct {
	StorageID    string
	VideoHash    string
	Timestamp    int64
	Chunks       []string
	ManifestPath string
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
		VideoHash:    result.VideoHash,
		Timestamp:    result.Timestamp,
		Chunks:       result.ChunkHashes,
		ManifestPath: result.ManifestPath,
	}, nil
}

func (c *StorageClient) GetChunk(hash string) ([]byte, error) {
	return c.processor.GetChunk(hash)
}

func (c *StorageClient) ComputeManifestHash(path string, authorPublicKey string) (string, error) {
	return videostorage.ComputeManifestHash(path, authorPublicKey)
}

func (c *StorageClient) SetManifestSignature(path string, authorPublicKey string, signature string) error {
	return videostorage.SetManifestSignature(path, authorPublicKey, signature)
}

func (c *StorageClient) SetManifestProof(path string, authorPublicKey string, signature string, videoHash string, timestamp int64) error {
	return videostorage.SetManifestProof(path, authorPublicKey, signature, videoHash, timestamp)
}
