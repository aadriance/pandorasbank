package quickstart

import (
	"os"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type QuickStartUI interface {

}

type quickStartUI struct {
	quickStarter QuickStarter
	knownPath string
	updateCustomApplicationPath func(string)

	window fyne.Window

	container *fyne.Container
	startBox *fyne.Container
	stopBox *fyne.Container
}

func NewQuickStartUI(myWindow fyne.Window, path string, updateCustomApplicationPath func(string)) (QuickStartUI, *fyne.Container) {
	ui := &quickStartUI{
		quickStarter: NewQuickStarter(),
		knownPath: path,
		window: myWindow,
		updateCustomApplicationPath: updateCustomApplicationPath,
		container: container.NewStack(),
	}

	ui.startBox = container.NewHBox(
		widget.NewButton("Start Game", func() {
			startGameButton(ui)
		}),
		widget.NewButton("Change Exe", func() {
			showFilePicker(ui, func() {})
		}),
	)

	// Note - currently, stop button replaced with status message
	ui.stopBox = container.NewHBox(
		widget.NewLabel("Game Running!"),
	)

	ui.container.Add(ui.startBox)

	return ui, ui.container
}

func startGameButton(ui *quickStartUI) {
	if !isValidPath(ui.knownPath) {
		confirmFilePicker(ui)
		return
	}

	err := ui.quickStarter.Start(ui.knownPath, func() {
		stopGameButton(ui)
	})

	if err != nil {
		// Failure to run?  Bad path maybe
		ui.knownPath = ""
		confirmFilePicker(ui)
		// Good luck to whomever forces this recursion enough times to stackoverflow
		return
	}

	ui.container.RemoveAll()
	ui.container.Add(ui.stopBox)
	ui.container.Refresh()
}

func confirmFilePicker(ui *quickStartUI) { 
	dialog.ShowConfirm("Game executable not located or invalid", "Could not locate your Elestrals Awakened Playtest by default.\n\nWant to select the Elestrals Awakened Playtest executable on your system?", func(result bool) {
		if result {
			showFilePicker(ui, func() {
				startGameButton(ui)
			})
		}
	}, ui.window)
}

func showFilePicker(ui *quickStartUI, onSuccess func()) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if reader == nil {
			// No file selected
			return
		}
		defer reader.Close()

		if err != nil {
			dialog.ShowError(err, ui.window)
			return
		}

		ui.knownPath = reader.URI().Path()
		ui.updateCustomApplicationPath(ui.knownPath)

		onSuccess()
	}, ui.window)
}

func stopGameButton(ui *quickStartUI) {
	ui.quickStarter.Stop()

	ui.container.RemoveAll()
	ui.container.Add(ui.startBox)

	// Possible call by goroutine
	fyne.Do(func() { ui.container.Refresh() })
}


func isValidPath(path string) bool {
  if path == "" {
		return false
	}

	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
