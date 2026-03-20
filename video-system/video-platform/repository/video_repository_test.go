package repository

import (
	"path/filepath"
	"testing"

	"video-platform/db"
	"video-platform/model"
)

func TestVideoRepositorySaveAndFind(t *testing.T) {
	tmp := t.TempDir()
	database, err := db.Open(filepath.Join(tmp, "platform.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer database.Close()
	repo := NewVideoRepository(database)

	v := model.Video{
		ID:       "v1",
		Title:    "hello",
		Filename: "hello.mp4",
	}
	if err := repo.Save(v); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, ok := repo.FindByID("v1")
	if !ok {
		t.Fatalf("expected video to exist")
	}
	if got.Title != "hello" {
		t.Fatalf("expected title hello, got %s", got.Title)
	}

	all := repo.FindAll()
	if len(all) != 1 {
		t.Fatalf("expected 1 video, got %d", len(all))
	}
}
