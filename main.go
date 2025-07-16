package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"time"
	"unicode"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/oto"
	"github.com/qeesung/image2ascii/convert"
	"github.com/tosone/minimp3"
)

func greeter() {
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	ascII("         Y2Go ")
	fmt.Print("V0.2 Dev")
	fmt.Println("\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Print("\nTo play music - enter a full path to an .mp3 file:")
}

// /home/daniil/Downloads/OldFlavours.mp3

func main() {
	greeter()
	path := scanner() // Bufio scan

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

func scanner() string {
	reader := bufio.NewScanner(os.Stdin)
	_ = reader.Scan()
	path := reader.Text()
	fmt.Print("\n")
	return path
}

func songLength(dec *minimp3.Decoder, data []byte) float64 {
	bytesPerSample := 2
	channels := dec.Channels
	sampleRate := dec.SampleRate

	totalSamples := len(data) / (bytesPerSample * channels)
	durationSeconds := float64(totalSamples) / float64(sampleRate)

	return durationSeconds
}

//func lenghtConverter(length float64) {
//	var minsfloat float64
//	minsfloat = length / 60.0
//	mins := int(minsfloat)
//	secsfloat := length - (float64(mins) * 60)
//	fmt.Printf("the length is %v mins %.2f secs\n", mins, secsfloat)
//}

func printMetadata(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	options := convert.DefaultOptions
	options.FixedWidth = 96
	options.FixedHeight = 30

	metadata, err := tag.ReadFrom(file)
	if err != nil {
		log.Fatal(err)
	}
	pic := metadata.Picture()
	if pic != nil {
		img, _ := decodeImage(pic.Data, pic.MIMEType)
		asc := convert.NewImageConverter()
		fmt.Print(asc.Image2ASCIIString(img, &options))
	}
	ascII(filterSpecialSymbols(fmt.Sprintf(metadata.Title())))
	ascIIlow(filterSpecialSymbols(fmt.Sprintf("By %v", metadata.Artist())))
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

func ascII(text string) {
	myFigure := figure.NewFigure(text, "nancyj", true)
	myFigure.Print()
}

func ascIIlow(text string) {
	myFigure := figure.NewFigure(text, "italic", true)
	myFigure.Print()
}

func filterSpecialSymbols(input string) string {
	var result []rune
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '(' || r == ')' || r == '_' || r == '-' {
			result = append(result, r)
		}
	}
	return string(result)
}
