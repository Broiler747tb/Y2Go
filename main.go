package main

import (
	"Y2Go/Player"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowden/tag"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	a := app.NewWithID("com.Player.Y2Go")
	w := a.NewWindow("Y2Go")

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter a file path:")

	Data := make(chan tag.Metadata, 3)
	Position := make(chan float64, 10)
	SetPosition := make(chan float64, 5)
	Stop := make(chan bool, 3)

	f, _ := os.Open("Assets/defaultCoverArt.jpg")
	ima, _, _ := image.Decode(f)
	f.Close()

	albumCover := canvas.NewImageFromImage(nil)
	albumCover.Image = ima
	albumCover.FillMode = canvas.ImageFillContain
	albumCover.SetMinSize(fyne.NewSize(100, 100))

	coverContainer := container.NewMax(albumCover)

	var userInput string
	button := widget.NewButton("Add to the queue/Play!", func() {
		go Player.Play(userInput, Data, Position, SetPosition, Stop)
	})

	selectFile := widget.NewButton("Select a file:", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader == nil {
				fmt.Println("file selection canceled")
				return
			}
			defer reader.Close()

			userInput = reader.URI().Path()
			fmt.Println(userInput)
		}, w)
		fileDialog.Resize(fyne.NewSize(9999, 9999))
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".mp3", ".wav", ".flac", ".ogg"}))
		fileDialog.Show()
	})

	go func() {
		for meta := range Data {
			fmt.Println("Received metadata:", meta.Title())
			pic := meta.Picture()
			if pic == nil {
				fmt.Println("No album art in metadata")
				f, _ := os.Open("Assets/defaultCoverArt.jpg")
				ima, _, _ = image.Decode(f)
				albumCover.Image = ima
				f.Close()
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

	slider := widget.NewSlider(0.0, 200.0)

	var isInternalChange bool
	var userInputs bool

	slider.OnChanged = func(value float64) {
		userInputs = true
	}
	slider.OnChangeEnded = func(value float64) {
		userInputs = false
		if isInternalChange {
			return
		} else {
			select {
			case SetPosition <- value:
			default:
			}
		}
	}

	go func() {
		for prog := range Position {
			fyne.Do(func() {
				if !userInputs {
					isInternalChange = true
					slider.SetValue(prog * 200)
					isInternalChange = false
				}
			})
		}
	}()

	if desk, ok := a.(desktop.App); ok {
		iconBytes, err := os.ReadFile("Assets/trayIcon.png")
		if err == nil {
			iconRes := fyne.NewStaticResource("Assets/trayIcon.png", iconBytes)
			a.SetIcon(iconRes) // Set the icon on the app
		}
		m := fyne.NewMenu("Y2Go",
			fyne.NewMenuItem("Show", func() {
				w.Show()
			}))
		desk.SetSystemTrayMenu(m)
	}

	Play := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		Stop <- true
	})
	Play.Resize(fyne.Size{20, 20})
	Nb := container.NewHBox(layout.NewSpacer(), Play, layout.NewSpacer())
	SliderAndButtons := container.NewBorder(Nb, slider, nil, nil)

	addWindow := container.New(layout.NewBorderLayout(nil, button, nil, nil), selectFile, button)
	playWindow := container.New(layout.NewBorderLayout(layout.NewSpacer(), SliderAndButtons, nil, nil), albumCover, SliderAndButtons)

	tabs := container.NewAppTabs(
		container.NewTabItem("Player", playWindow),
		container.NewTabItem("Add songs:", addWindow),
		container.NewTabItem("Queue", widget.NewLabel("WIP")),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(400, 400))
	w.SetCloseIntercept(func() {
		w.Hide()
	})
	w.ShowAndRun()
}
