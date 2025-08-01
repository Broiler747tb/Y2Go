package Player

import "C"
import (
	"github.com/dhowden/tag"
	"github.com/eiannone/keyboard"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"

	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Play(path string, metadata chan tag.Metadata, Playended chan bool, position chan float64) {
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
	go func() {
		for {
			percent := float64(streamer.Position()) / float64(streamer.Len())
			time.Sleep(time.Millisecond * 250)
			position <- percent
		}
	}()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
		done <- true
	})))

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

func stopper() {
	keysEvents, _ := keyboard.GetKeys(10)
	defer func() {
		_ = keyboard.Close()
	}()
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
		if event.Key == keyboard.KeyEsc {
			break
		}
	}
}
