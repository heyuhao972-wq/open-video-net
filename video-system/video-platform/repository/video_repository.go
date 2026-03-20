package repository

import (
	"database/sql"
	"encoding/json"

	"video-platform/model"
)

type VideoRepository struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) Save(video model.Video) error {
	tags, _ := json.Marshal(video.Tags)
	chunks, _ := json.Marshal(video.Chunks)
	_, err := r.db.Exec(
		`INSERT INTO videos
		(video_id, platform_id, storage_id, title, description, filename, file_path, tags, author_id, author_public_key, author_signature, proof_timestamp, video_hash, chunks, manifest, manifest_hash, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(video_id) DO UPDATE SET
		  title=excluded.title,
		  description=excluded.description,
		  tags=excluded.tags,
		  author_id=excluded.author_id,
		  author_public_key=excluded.author_public_key,
		  author_signature=excluded.author_signature,
		  proof_timestamp=excluded.proof_timestamp,
		  video_hash=excluded.video_hash,
		  chunks=excluded.chunks,
		  manifest=excluded.manifest,
		  manifest_hash=excluded.manifest_hash,
		  platform_id=excluded.platform_id`,
		video.ID, video.PlatformID, video.StorageID, video.Title, video.Description, video.Filename, video.FilePath,
		string(tags), video.AuthorID, video.AuthorPublicKey, video.AuthorSignature, video.ProofTimestamp, video.VideoHash,
		string(chunks), video.Manifest, video.ManifestHash, video.CreatedAt,
	)
	return err
}

func (r *VideoRepository) FindAll() []model.Video {
	rows, err := r.db.Query(`SELECT video_id, platform_id, storage_id, title, description, filename, file_path, tags, author_id, author_public_key, author_signature, proof_timestamp, video_hash, chunks, manifest, manifest_hash, created_at FROM videos`)
	if err != nil {
		return []model.Video{}
	}
	defer rows.Close()

	out := []model.Video{}
	for rows.Next() {
		var v model.Video
		var tags string
		var chunks string
		_ = rows.Scan(&v.ID, &v.PlatformID, &v.StorageID, &v.Title, &v.Description, &v.Filename, &v.FilePath, &tags,
			&v.AuthorID, &v.AuthorPublicKey, &v.AuthorSignature, &v.ProofTimestamp, &v.VideoHash, &chunks, &v.Manifest, &v.ManifestHash, &v.CreatedAt)
		_ = json.Unmarshal([]byte(tags), &v.Tags)
		_ = json.Unmarshal([]byte(chunks), &v.Chunks)
		out = append(out, v)
	}
	return out
}

func (r *VideoRepository) FindByID(id string) (model.Video, bool) {
	var v model.Video
	var tags string
	var chunks string
	err := r.db.QueryRow(
		`SELECT video_id, platform_id, storage_id, title, description, filename, file_path, tags, author_id, author_public_key, author_signature, proof_timestamp, video_hash, chunks, manifest, manifest_hash, created_at FROM videos WHERE video_id = ?`,
		id,
	).Scan(&v.ID, &v.PlatformID, &v.StorageID, &v.Title, &v.Description, &v.Filename, &v.FilePath, &tags,
		&v.AuthorID, &v.AuthorPublicKey, &v.AuthorSignature, &v.ProofTimestamp, &v.VideoHash, &chunks, &v.Manifest, &v.ManifestHash, &v.CreatedAt)
	if err != nil {
		return model.Video{}, false
	}
	_ = json.Unmarshal([]byte(tags), &v.Tags)
	_ = json.Unmarshal([]byte(chunks), &v.Chunks)
	return v, true
}
