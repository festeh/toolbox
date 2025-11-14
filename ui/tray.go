package ui

import (
	_ "embed"

	"fyne.io/systray"
)

//go:embed icon.png
var iconData []byte

var (
	showChan chan bool
	quitChan chan bool
)

// InitTray initializes channels for communication
func InitTray() (chan bool, chan bool) {
	showChan = make(chan bool, 1)
	quitChan = make(chan bool, 1)
	return showChan, quitChan
}

// RunTray runs the system tray (blocking call)
func RunTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Toolbox")
	systray.SetTooltip("Toolbox")

	// Menu items
	mShow := systray.AddMenuItem("Show", "Show the window")
	mHide := systray.AddMenuItem("Hide", "Hide the window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				showChan <- true
			case <-mHide.ClickedCh:
				showChan <- false
			case <-mQuit.ClickedCh:
				quitChan <- true
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// Cleanup
}
