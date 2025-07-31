package main

import (
	"Y2Go/Cli"
	"Y2Go/Flac"
	"Y2Go/Mp3"
	"Y2Go/Vorbis"
	"Y2Go/Wav"
	"strings"
)

func main() {
	path := CLI.GreeterAndSelecter()
	var data *CLI.Metadata

	if strings.HasSuffix(path, ".mp3") {
		Mp3.PlayMp3(path)
	}
	if strings.HasSuffix(path, ".flac") {
		Flac.PlayFlac(path)
	}
	if strings.HasSuffix(path, ".wav") {
		Wav.PlayWav(path)
	}
	if strings.HasSuffix(path, ".ogg") {
		Vorbis.PlayVorbis(path)
	}
	CLI.PrintMetadata(*data)
}
