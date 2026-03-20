package p2p

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

type ChunkFetcher struct {
	platformBase string
	p2pBase      string
	cache        *chunkCache
	maxParallel  int
	client       *http.Client
	retry        int
}

type manifestResponse struct {
	VideoID         string `json:"video_id"`
	AuthorPublicKey string `json:"author_public_key"`
	Signature       string `json:"signature"`
	Chunks          []struct {
		Hash  string `json:"hash"`
		Index int    `json:"index"`
		Size  int    `json:"size"`
	} `json:"chunks"`
}

func NewChunkFetcher(platformBase string, p2pBase string, maxParallel int, cacheSize int, timeoutMs int, retry int) *ChunkFetcher {
	if maxParallel <= 0 {
		maxParallel = 6
	}
	if cacheSize <= 0 {
		cacheSize = 128
	}
	if timeoutMs <= 0 {
		timeoutMs = 4000
	}
	if retry < 0 {
		retry = 0
	}

	return &ChunkFetcher{
		platformBase: platformBase,
		p2pBase:      p2pBase,
		cache:        newChunkCache(cacheSize),
		maxParallel:  maxParallel,
		client: &http.Client{
			Timeout: time.Duration(timeoutMs) * time.Millisecond,
		},
		retry: retry,
	}
}

func (f *ChunkFetcher) FetchChunks(videoId string) ([][]byte, error) {
	if f.platformBase == "" {
		return nil, errors.New("platform base missing")
	}

	manifestURL := fmt.Sprintf("%s/video/%s/manifest", f.platformBase, videoId)
	res, err := f.client.Get(manifestURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest status: %d", res.StatusCode)
	}

	var m manifestResponse
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, err
	}

	if err := verifyManifestSignature(m); err != nil {
		return nil, err
	}

	if len(m.Chunks) == 0 {
		return nil, errors.New("no chunks")
	}

	sort.Slice(m.Chunks, func(i, j int) bool {
		return m.Chunks[i].Index < m.Chunks[j].Index
	})

	chunks := make([][]byte, len(m.Chunks))
	sem := make(chan struct{}, f.maxParallel)
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for i, c := range m.Chunks {
		wg.Add(1)
		go func(idx int, hash string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if data, ok := f.cache.Get(hash); ok {
				chunks[idx] = data
				return
			}

			data, err := f.fetchChunk(hash)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
			f.cache.Set(hash, data)
			chunks[idx] = data
		}(i, c.Hash)
	}

	wg.Wait()
	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	return chunks, nil
}

func verifyManifestSignature(m manifestResponse) error {
	if m.AuthorPublicKey == "" || m.Signature == "" {
		return errors.New("manifest signature missing")
	}

	pubBytes, err := base64.StdEncoding.DecodeString(m.AuthorPublicKey)
	if err != nil {
		return fmt.Errorf("invalid author public key: %w", err)
	}

	sigBytes, err := base64.StdEncoding.DecodeString(m.Signature)
	if err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}

	hashHex, err := computeManifestHash(m)
	if err != nil {
		return err
	}
	hashBytes, err := hex.DecodeString(hashHex)
	if err != nil {
		return err
	}

	if !ed25519.Verify(ed25519.PublicKey(pubBytes), hashBytes, sigBytes) {
		return errors.New("manifest signature verification failed")
	}

	return nil
}

func computeManifestHash(m manifestResponse) (string, error) {
	tmp := manifestResponse{
		VideoID:         m.VideoID,
		AuthorPublicKey: m.AuthorPublicKey,
		Signature:       "",
		Chunks:          m.Chunks,
	}

	data, err := json.Marshal(tmp)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func (f *ChunkFetcher) fetchChunk(hash string) ([]byte, error) {
	if f.p2pBase != "" {
		if data, ok := f.fetchFromGateway(hash); ok {
			return data, nil
		}
	}
	chunkURL := fmt.Sprintf("%s/chunk/%s", f.platformBase, hash)
	var lastErr error
	for attempt := 0; attempt <= f.retry; attempt++ {
		chunkRes, err := f.client.Get(chunkURL)
		if err != nil {
			lastErr = err
			continue
		}
		if chunkRes.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("chunk status: %d", chunkRes.StatusCode)
			chunkRes.Body.Close()
			continue
		}
		data, err := io.ReadAll(chunkRes.Body)
		chunkRes.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		return data, nil
	}
	return nil, lastErr
}

func (f *ChunkFetcher) fetchFromGateway(hash string) ([]byte, bool) {
	p2pURL := fmt.Sprintf("%s/chunk/%s", f.p2pBase, hash)
	resp, err := f.client.Get(p2pURL)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, false
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false
	}
	return data, true
}

type chunkCache struct {
	capacity int
	order    []string
	items    map[string][]byte
	lock     sync.Mutex
}

func newChunkCache(capacity int) *chunkCache {
	return &chunkCache{
		capacity: capacity,
		order:    make([]string, 0, capacity),
		items:    make(map[string][]byte),
	}
}

func (c *chunkCache) Get(key string) ([]byte, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	v, ok := c.items[key]
	return v, ok
}

func (c *chunkCache) Set(key string, value []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.items[key]; ok {
		c.items[key] = value
		return
	}

	if len(c.order) >= c.capacity {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.items, oldest)
	}

	c.items[key] = value
	c.order = append(c.order, key)
}
