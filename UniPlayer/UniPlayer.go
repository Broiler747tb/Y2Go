package UniPlayer

import (
	CLI "Y2Go/Cli"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Play(path string) {
	ext := filepath.Ext(path)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	var Meta CLI.Metadata

	m, err := tag.ReadFrom(f)
	if err != nil {
		log.Fatal(err)
	}
	if m != nil {
		Meta = CLI.Metadata{
			Song:   m.Title(),
			Artist: m.Artist(),
			Album:  m.Album(),
			Year:   m.Year(),
			Genre:  m.Genre(),
		}
		if pic := m.Picture(); pic != nil {
			Meta.Picture = *pic
		}
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format := extensionSwitcher(f, ext)
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
		done <- true
	})))
	CLI.PrintMetadata(&Meta)
	<-done
}

func extensionSwitcher(f io.ReadCloser, ext string) (beep.StreamSeekCloser, beep.Format) {
	switch ext {
	case ".mp3":
		streamer, format, err := mp3.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		return streamer, format
	case ".flac":
		streamer, format, err := flac.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		return streamer, format
	case ".wav":
		streamer, format, err := wav.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		return streamer, format
	case ".ogg":
		streamer, format, err := vorbis.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		return streamer, format
	default:
		fmt.Println("Error: Not supported file format!")
		var streamer beep.StreamSeekCloser
		var format beep.Format
		return streamer, format
	}
}
