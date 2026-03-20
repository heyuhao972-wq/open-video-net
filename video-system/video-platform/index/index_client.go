package index

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"video-platform/model"
)

type Client struct {
	base string
}

func NewClient(base string) *Client {
	return &Client{base: base}
}

func (c *Client) UpsertVideo(v model.Video) error {
	body, err := json.Marshal(map[string]interface{}{
		"id":                v.ID,
		"platform_id":       v.PlatformID,
		"title":             v.Title,
		"description":       v.Description,
		"tags":              v.Tags,
		"author_id":         v.AuthorID,
		"author_public_key": v.AuthorPublicKey,
		"manifest_hash":     v.ManifestHash,
		"views":             0,
		"created_at":        v.CreatedAt,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(c.base+"/video", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("index status: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) IncrementViews(id string) error {
	body, err := json.Marshal(map[string]interface{}{
		"id":    id,
		"views": 1,
		"op":    "inc",
		"field": "views",
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(c.base+"/video/views", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("index status: %d", resp.StatusCode)
	}
	return nil
}
