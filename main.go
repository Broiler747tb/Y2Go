package main

import (
	"bytes"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"time"
)

// /home/daniil/Downloads/OldFlavours.mp3

func main() {
	path := "/home/daniil/Downloads/OldFlavours.mp3"

	var err error
	var file []byte
	if file, err = os.ReadFile(path); err != nil {
		log.Fatal(err)
	} // Scanning the file for minimp3

	printMetadata(path) // Reading the metadata from a file and printing it

	var dec *minimp3.Decoder
	var data []byte
	if dec, data, err = minimp3.DecodeFull(file); err != nil {
		log.Fatal(err)
	} // Decoding the scanned mp3 file

	//leng := songLength(dec, data)
	//go lenghtConverter(leng)

	var context *oto.Context
	if context, err = oto.NewContext(dec.SampleRate, dec.Channels, 2, 256); err != nil {
		log.Fatal(err)
	} // I honestly dont know what is ts

	var player = context.NewPlayer()
	player.Write(data) // Playing the decoded thing

	<-time.After(time.Second)

	dec.Close() // Player ending and after 1 second closing the decoder
	if err = player.Close(); err != nil {
		log.Fatal(err)
	}
}

func printMetadata(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	metadata, err := tag.ReadFrom(file)
	if err != nil {
		log.Fatal(err)
	}
	pic := metadata.Picture()
	if pic != nil {
		decodeImage(pic.Data, pic.MIMEType)
	}
	fmt.Sprintf(metadata.Title())
	fmt.Sprintf("By %v", metadata.Artist())
}

func decodeImage(data []byte, mimeType string) (image.Image, error) {
	reader := bytes.NewReader(data)
	switch mimeType {
	case "image/jpeg", "image/jpg":
		return jpeg.Decode(reader)
	case "image/png":
		return png.Decode(reader)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", mimeType)
	}
}
