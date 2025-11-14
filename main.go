package main

import (
	_ "embed"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed style.css
var styleCSS string

func main() {
	app := gtk.NewApplication("com.example.toolbox", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

type TmuxWindow struct {
	Session  string
	Index    string
	Name     string
	Activity string
}

func activate(app *gtk.Application) {
	// Load Catppuccin Mocha CSS theme
	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		loadCSS(styleCSS),
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	// Get tmux windows list
	windows := getTmuxWindows()

	// Create a box for the table
	box := gtk.NewBox(gtk.OrientationVertical, 0)

	// Create header
	headerBox := gtk.NewBox(gtk.OrientationHorizontal, 0)
	headerBox.AddCSSClass("table-header")
	headerBox.SetHomogeneous(false)

	sessionHeader := gtk.NewLabel("Session")
	sessionHeader.SetXAlign(0)
	sessionHeader.SetSizeRequest(150, -1)
	sessionHeader.AddCSSClass("header-cell")

	indexHeader := gtk.NewLabel("Index")
	indexHeader.SetXAlign(0)
	indexHeader.SetSizeRequest(80, -1)
	indexHeader.AddCSSClass("header-cell")

	nameHeader := gtk.NewLabel("Window Name")
	nameHeader.SetXAlign(0)
	nameHeader.SetHExpand(true)
	nameHeader.AddCSSClass("header-cell")

	activityHeader := gtk.NewLabel("Activity")
	activityHeader.SetXAlign(0)
	activityHeader.SetSizeRequest(200, -1)
	activityHeader.AddCSSClass("header-cell")

	headerBox.Append(sessionHeader)
	headerBox.Append(indexHeader)
	headerBox.Append(nameHeader)
	headerBox.Append(activityHeader)

	box.Append(headerBox)

	// Create scrollable content area
	contentBox := gtk.NewBox(gtk.OrientationVertical, 0)
	contentBox.AddCSSClass("table-content")

	// Add rows
	for _, window := range windows {
		rowBox := gtk.NewBox(gtk.OrientationHorizontal, 0)
		rowBox.AddCSSClass("table-row")
		rowBox.SetHomogeneous(false)

		sessionLabel := gtk.NewLabel(window.Session)
		sessionLabel.SetXAlign(0)
		sessionLabel.SetSizeRequest(150, -1)
		sessionLabel.AddCSSClass("table-cell")

		indexLabel := gtk.NewLabel(window.Index)
		indexLabel.SetXAlign(0)
		indexLabel.SetSizeRequest(80, -1)
		indexLabel.AddCSSClass("table-cell")

		nameLabel := gtk.NewLabel(window.Name)
		nameLabel.SetXAlign(0)
		nameLabel.SetHExpand(true)
		nameLabel.AddCSSClass("table-cell")

		activityLabel := gtk.NewLabel(formatTimestamp(window.Activity))
		activityLabel.SetXAlign(0)
		activityLabel.SetSizeRequest(200, -1)
		activityLabel.AddCSSClass("table-cell")

		rowBox.Append(sessionLabel)
		rowBox.Append(indexLabel)
		rowBox.Append(nameLabel)
		rowBox.Append(activityLabel)

		contentBox.Append(rowBox)
	}

	// Create scrolled window
	scrolledWindow := gtk.NewScrolledWindow()
	scrolledWindow.SetChild(contentBox)
	scrolledWindow.SetVExpand(true)
	scrolledWindow.SetHExpand(true)

	box.Append(scrolledWindow)

	// Create window
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Toolbox")
	window.SetDefaultSize(900, 600)
	window.SetChild(box)
	window.Show()
}

func getTmuxWindows() []TmuxWindow {
	cmd := exec.Command("tmux", "list-windows", "-a", "-F", "#{session_name}\t#{window_index}\t#{window_name}\t#{window_activity}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []TmuxWindow{{
			Session:  "Error",
			Index:    "",
			Name:     err.Error(),
			Activity: "",
		}}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	windows := make([]TmuxWindow, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 4 {
			windows = append(windows, TmuxWindow{
				Session:  parts[0],
				Index:    parts[1],
				Name:     parts[2],
				Activity: parts[3],
			})
		}
	}

	// Sort by activity timestamp (most recent first)
	sort.Slice(windows, func(i, j int) bool {
		activityI, _ := strconv.ParseInt(windows[i].Activity, 10, 64)
		activityJ, _ := strconv.ParseInt(windows[j].Activity, 10, 64)
		return activityI > activityJ
	})

	return windows
}

func formatTimestamp(ts string) string {
	timestamp, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return ts
	}

	t := time.Unix(timestamp, 0)

	// Format as "15:04 2 Jan"
	return t.Format("15:04 2 Jan")
}

func loadCSS(content string) *gtk.CSSProvider {
	prov := gtk.NewCSSProvider()
	prov.ConnectParsingError(func(sec *gtk.CSSSection, err error) {
		loc := sec.StartLocation()
		lines := strings.Split(content, "\n")
		log.Printf("CSS error (%v) at line: %q", err, lines[loc.Lines()])
	})
	prov.LoadFromData(content)
	return prov
}
