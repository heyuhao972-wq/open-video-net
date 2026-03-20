package main

import (
	"fmt"
	"os"

	"video-storage/internal/chunk"
	"video-storage/internal/manifest"
	"video-storage/internal/storage"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <video_file>")
		return
	}

	videoPath := os.Args[1]

	fmt.Println("Reading video:", videoPath)

	chunks, err := chunk.SplitFile(videoPath, 1024*1024) // 1MB chunk
	if err != nil {
		panic(err)
	}

	fmt.Println("Chunks created:", len(chunks))

	store, err := storage.NewLocalStore("./data/chunks")
	if err != nil {
		panic(err)
	}

	for _, c := range chunks {
		err := store.Save(c)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Chunks stored")

	manifest, err := manifest.BuildManifest(chunks, "hash", 0)
	if err != nil {
		panic(err)
	}

	err = manifest.Save("./data/manifest.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("Manifest generated")
}
