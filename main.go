package main

import (
	"Y2Go/Cli"
	"Y2Go/Flac"
	"Y2Go/Mp3"
	"Y2Go/Vorbis"
	"Y2Go/Wav"
	"fmt"
	"path/filepath"
)

func main() {
	path := CLI.GreeterAndSelecter()
	ext := filepath.Ext(path)
	switch ext {
	case ".mp3":
		Mp3.PlayMp3(path)
	case ".flac":
		Flac.PlayFlac(path)
	case ".wav":
		Wav.PlayWav(path)
	case ".ogg":
		Vorbis.PlayVorbis(path)
	default:
		fmt.Println("Error: Not supported file format!")
	}
}
