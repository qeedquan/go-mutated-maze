package main

import (
	"log"
	"path/filepath"

	"github.com/qeedquan/go-media/sdl/sdlmixer"
)

var (
	sounds = make(map[string]*sdlmixer.Chunk)
)

func loadSound(name string) *sdlmixer.Chunk {
	filename := filepath.Join(*dataDir, name)
	if chunk, found := sounds[filename]; found {
		return chunk
	}

	chunk, err := sdlmixer.LoadWAV(filename)
	if err != nil {
		log.Fatal(err)
	}
	sounds[filename] = chunk
	return chunk
}

func playSound(chunk *sdlmixer.Chunk) {
	if !*sfx || chunk == nil {
		return
	}
	chunk.PlayChannel(-1, 0)
}