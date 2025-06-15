package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	window      fyne.Window
	statusLabel *widget.Label
	dropArea    *widget.Label
}

func main() {
	myApp := app.New()

	w := myApp.NewWindow("Tesla group invoices")
	w.Resize(fyne.NewSize(600, 250))
	w.CenterOnScreen()

	appWindow := &App{
		window: w,
	}

	appWindow.makeUI()
	w.ShowAndRun()
}

func (a *App) makeUI() {
	instructions := widget.NewLabel(
		"Glisser / déposer le zip.\n" +
			"L'application va extraire tous les PDF dans un seul dossier.")
	instructions.Wrapping = fyne.TextWrapWord
	instructions.Alignment = fyne.TextAlignCenter

	a.dropArea = widget.NewLabel("Glisser le dossier zip")
	a.dropArea.Alignment = fyne.TextAlignCenter
	a.dropArea.TextStyle = fyne.TextStyle{Bold: true}

	selectButton := widget.NewButton("Ou selectionner un dossier zip", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
		}, a.window)
	})

	content := container.NewVBox(
		widget.NewSeparator(),
		instructions,
		widget.NewSeparator(),
		container.NewPadded(a.dropArea),
		widget.NewSeparator(),
		selectButton,
	)

	a.window.SetContent(container.NewPadded(content))
}
