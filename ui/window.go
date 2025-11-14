package ui

import (
	"log"
	"toolbox/tmux"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// CreateMainWindow creates and shows the main application window
func CreateMainWindow(app *gtk.Application) {
	// Load CSS
	LoadCSS()

	// Get tmux windows
	windows := tmux.GetWindows()

	// Create main container
	mainBox := gtk.NewBox(gtk.OrientationVertical, 0)

	// Create table view
	tableView := CreateTableView(windows)

	// Create search input
	searchInput := CreateSearchInput(tableView, windows)

	// Add widgets to main container
	mainBox.Append(searchInput.Entry)
	mainBox.Append(tableView.Container)

	// Create window
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Toolbox")
	window.SetDefaultSize(900, 600)
	window.SetChild(mainBox)

	// Add keyboard event handler with capture phase to intercept before search input
	evtController := gtk.NewEventControllerKey()
	evtController.SetPropagationPhase(gtk.PhaseCapture)
	evtController.ConnectKeyPressed(func(keyval, keycode uint, state gdk.ModifierType) bool {
		// Check if a number key 1-9 is pressed
		if keyval >= gdk.KEY_1 && keyval <= gdk.KEY_9 {
			id := int(keyval - gdk.KEY_1 + 1)
			win := tableView.GetWindowByID(id)
			if win != nil {
				log.Printf("Switching to window: %s:%s (%s)", win.Session, win.Index, win.Name)
				err := tmux.SwitchToWindow(win.Session, win.Index)
				if err != nil {
					log.Printf("Error switching to window: %v", err)
				}
			}
			return true // Stop event propagation
		}
		return false // Allow other keys to propagate
	})
	window.AddController(evtController)

	window.Show()
}

// Run starts the GTK application
func Run() int {
	app := gtk.NewApplication("com.example.toolbox", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { CreateMainWindow(app) })
	return app.Run(nil)
}
