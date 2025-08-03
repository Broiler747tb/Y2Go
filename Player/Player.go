package Player

import "C"
import (
	"github.com/dhowden/tag"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"

	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var speakerInitialized bool = false
var GlobalPlayEnded bool = true

func Play(path string, metadata chan tag.Metadata, position chan float64, setPosition chan float64) {
	for !GlobalPlayEnded {
		time.Sleep(time.Millisecond * 100)
	}
	var playEnded bool = true
	ext := filepath.Ext(path)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	m, err := tag.ReadFrom(f)
	if err != nil {
		log.Fatal(err)
	}
	metadata <- m
	_, err = f.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format := extensionSwitcher(f, ext)
	defer streamer.Close()

	if !speakerInitialized {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		speakerInitialized = true
	}
	playEnded = false
	GlobalPlayEnded = false

	done := make(chan bool)
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
		done <- true
	})))

	var seeking bool

	go func() {
		for set := range setPosition {
			if !playEnded {
				seeking = true
				pos := int(float64(streamer.Len()) * (set / 200.0))
				streamer.Seek(pos)
				seeking = false
			}
		}
	}()

	go func() {
		for {
			if !seeking {
				if !playEnded {
					time.Sleep(time.Millisecond * 200)
					percent := float64(streamer.Position()) / float64(streamer.Len())
					select {
					case position <- percent:
					default:
					}
				}
			} else {
				break
			}
		}
	}()

	<-done
	playEnded = true
	GlobalPlayEnded = true
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
