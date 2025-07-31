package CLI

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"unicode"

	"github.com/common-nighthawk/go-figure"
	"github.com/dhowden/tag"
)

type Metadata struct {
	Song    string
	Artist  string
	Album   string
	Year    int
	Genre   string
	Picture tag.Picture
}

func GreeterAndSelecter() string {
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	ascII("         Y2Go")
	fmt.Printf("\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ V0.5.1 Dev\n")
	path := Promt("Enter a full path to a music file:")
	return path
}

func Promt(promt string) string {
	fmt.Print(promt)
	scanner := bufio.NewScanner(os.Stdin)
	_ = scanner.Scan()
	input := scanner.Text()
	return input
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
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '(' || r == ')' || r == '_' || r == '-' || r == '?' || r == '!' || r == '@' || r == '#' || r == '&' || r == '*' || r == '"' || r == '/' || r == '|' {
			result = append(result, r)
		}
	}
	return string(result)
}

func PrintMetadata(m *Metadata) {
	//if m.Picture != nil {
	//decodeImage(m.Picture.Data, m.Picture.MIMEType)
	//}
	ascII(filterSpecialSymbols(fmt.Sprintf(m.Song)))
	ascIIlow(filterSpecialSymbols(fmt.Sprintf("By %v", m.Artist)))
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
