package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	window   fyne.Window
	dropArea *widget.Label
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
			"L'application va extraire tous les PDF dans un seul dossier sur le bureau.")
	instructions.Wrapping = fyne.TextWrapWord
	instructions.Alignment = fyne.TextAlignCenter

	a.dropArea = widget.NewLabel("Glisser le dossier zip")
	a.dropArea.Alignment = fyne.TextAlignCenter
	a.dropArea.TextStyle = fyne.TextStyle{Bold: true}

	selectButton := widget.NewButton("Ou selectionner un dossier zip", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				log.Printf("Erreur lors de l'ouverture du fichier : %v", err)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			zipFilePath := reader.URI().Path()
			a.handleZipFile(zipFilePath)
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

func (a *App) handleZipFile(zipFilePath string) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		a.dropArea.SetText(fmt.Sprintf("Impossible de trouver le répertoire de l'utilisateur: %v", err))
		return
	}

	desktopPath := filepath.Join(userHomeDir, "Desktop")
	currentDate := time.Now().Format("20060102150405")
	folderName := fmt.Sprintf("tesla_facture_%s", currentDate)

	outputDir := filepath.Join(desktopPath, folderName)
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		a.dropArea.SetText(fmt.Sprintf("Impossible de créer le dossier'%s': %v", outputDir, err))
		return
	}

	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		a.dropArea.SetText(fmt.Sprintf("Impossible d'ouvrir le zip '%s': %v", filepath.Base(zipFilePath), err))
		return
	}
	defer r.Close()

	pdfCount := 0
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		if strings.HasSuffix(strings.ToLower(f.Name), ".pdf") {
			rc, err := f.Open()
			if err != nil {
				a.dropArea.SetText(fmt.Sprintf("Erreur lors de la lecture des pdf '%s': %v", f.Name, err))
				continue
			}

			outputPath := filepath.Join(outputDir, filepath.Base(f.Name))

			outFile, err := os.Create(outputPath)
			if err != nil {
				a.dropArea.SetText(fmt.Sprintf("Erreur lors de la création du fichier '%s': %v", outputPath, err))
				rc.Close()
				continue
			}

			_, err = io.Copy(outFile, rc)
			if err != nil {
				a.dropArea.SetText(fmt.Sprintf("Erreur lors de la copie du fichier '%s' vers '%s': %v", f.Name, outputPath, err))
			} else {
				pdfCount++
				a.dropArea.SetText(fmt.Sprintf("Extraction en cours... %d PDF(s) extraits.", pdfCount))
			}

			outFile.Close()
			rc.Close()
		}
	}

	if pdfCount > 0 {
		a.dropArea.SetText(fmt.Sprintf("%d PDF extrait: '%s'!", pdfCount, outputDir))
	} else {
		a.dropArea.SetText("Erreur lors de l'extraction")
	}
}
