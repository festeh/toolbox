package main

import (
	"os"
	"toolbox/ui"
)

func main() {
	// Initialize tray communication channels
	showChan, quitChan := ui.InitTray()

	// Run system tray in a goroutine
	go ui.RunTray()

	// Run GTK application (blocking)
	os.Exit(ui.Run(showChan, quitChan))
}
