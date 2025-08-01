package Utils

import (
	"bytes"
	"fmt"
	"github.com/dhowden/tag"
	"image"
	"log"
	"unicode"
)

type Metadata struct {
	Song    string
	Artist  string
	Album   string
	Year    int
	Genre   string
	Picture tag.Picture
}

func FilterSpecialSymbols(input string) string {
	var result []rune
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '(' || r == ')' || r == '_' || r == '-' || r == '?' || r == '!' || r == '@' || r == '#' || r == '&' || r == '*' || r == '"' || r == '/' || r == '|' {
			result = append(result, r)
		} else {
			fmt.Println("Error: Unsupported symbol in the Name/Metadata")
		}
	}
	return string(result)
}

func ImageFromMetadata(pic tag.Picture) image.Image {
	img, _, err := image.Decode(bytes.NewReader(pic.Data))
	if err != nil {
		log.Fatal("Error decoding image:", err)
	}

	return img
}
