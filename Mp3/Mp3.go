package Mp3

import (
	CLI "Y2Go/Cli"
	"github.com/dhowden/tag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
	"time"
)

func PlayMp3(path string) {
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

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(ctrl)

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	CLI.PrintMetadata(Meta)
	<-done
}
