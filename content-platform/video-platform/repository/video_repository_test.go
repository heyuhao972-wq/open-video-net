package repository

import (
	"testing"

	"video-platform/model"
)

func TestVideoRepositorySaveAndFind(t *testing.T) {
	repo := NewVideoRepository()

	v := model.Video{
		ID:       "v1",
		Title:    "hello",
		Filename: "hello.mp4",
	}
	repo.Save(v)

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
