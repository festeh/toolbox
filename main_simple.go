package main

import (
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	app := gtk.NewApplication("com.example.helloworld.simple", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activateSimple(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activateSimple(app *gtk.Application) {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Hello World Test")
	window.SetChild(gtk.NewLabel("Hello World"))
	window.SetDefaultSize(700, 500)
	window.Show()
}
