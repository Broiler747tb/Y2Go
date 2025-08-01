package main

import (
	"Y2Go/Player"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowden/tag"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	a := app.New()
	w := a.NewWindow("Y2Go")

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter a file path:")

	Data := make(chan tag.Metadata, 3)
	Done := make(chan bool, 3)
	Position := make(chan float64, 10)

	f, _ := os.Open("4a51f2bcb67da9e5f941ffcc89f9bf00.jpg")
	ima, _, _ := image.Decode(f)
	f.Close()

	albumCover := canvas.NewImageFromImage(nil)
	albumCover.Image = ima
	albumCover.FillMode = canvas.ImageFillContain
	albumCover.SetMinSize(fyne.NewSize(100, 100))
	coverContainer := container.NewMax(albumCover)

	var userInput string
	button := widget.NewButton("Add to the queue/Play!", func() {
		userInput = entry.Text
		fmt.Println("Button pressed. Playing:", userInput)
		go Player.Play(userInput, Data, Done, Position)
	})

	progress := widget.NewProgressBar()

	// Listen for metadata and update UI
	go func() {
		for meta := range Data {
			fmt.Println("Received metadata:", meta.Title())
			pic := meta.Picture()
			if pic == nil {
				fmt.Println("No album art in metadata")
				continue
			}

			img, _, err := image.Decode(bytes.NewReader(pic.Data))
			if err != nil {
				fmt.Println("Error decoding image:", err)
				continue
			}

			fyne.Do(func() {
				fmt.Println("Updating image")
				albumCover.Image = img
				albumCover.Refresh()
				coverContainer.Refresh()
				cont := fmt.Sprint(`"`, meta.Title(), `"`, " by ", meta.Artist())
				a.SendNotification(&fyne.Notification{
					Title:   "Now Playing:",
					Content: cont,
				})
			})
		}
	}()

	go func() {
		for prog := range Position {
			fyne.Do(func() {
				progress.SetValue(prog)
			})
		}
	}()

	entryContainer := container.New(layout.NewMaxLayout(), entry)
	addWindow := container.New(layout.NewBorderLayout(nil, button, nil, nil), entryContainer, button)
	playWindow := container.New(layout.NewBorderLayout(nil, progress, nil, nil), coverContainer, progress)

	tabs := container.NewAppTabs(
		container.NewTabItem("Player", playWindow),
		container.NewTabItem("Add songs:", addWindow),
		container.NewTabItem("Queue", widget.NewLabel("WIP")),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(300, 300))
	w.ShowAndRun()
}

// /home/daniil/Downloads/OldFlavours.mp3
