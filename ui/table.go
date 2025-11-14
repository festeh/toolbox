package ui

import (
	"toolbox/tmux"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type TableView struct {
	Container      *gtk.Box
	ContentBox     *gtk.Box
	allWindows     []tmux.TmuxWindow
	displayedWindows []tmux.TmuxWindow
}

// CreateTableView creates a table view with header and scrollable content
func CreateTableView(windows []tmux.TmuxWindow) *TableView {
	tv := &TableView{
		Container:  gtk.NewBox(gtk.OrientationVertical, 0),
		allWindows: windows,
	}

	// Create header
	headerBox := gtk.NewBox(gtk.OrientationHorizontal, 0)
	headerBox.AddCSSClass("table-header")
	headerBox.SetHomogeneous(false)

	idHeader := gtk.NewLabel("#")
	idHeader.SetXAlign(0)
	idHeader.SetSizeRequest(40, -1)
	idHeader.AddCSSClass("header-cell")

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

	headerBox.Append(idHeader)
	headerBox.Append(sessionHeader)
	headerBox.Append(indexHeader)
	headerBox.Append(nameHeader)
	headerBox.Append(activityHeader)

	tv.Container.Append(headerBox)

	// Create scrollable content area
	tv.ContentBox = gtk.NewBox(gtk.OrientationVertical, 0)
	tv.ContentBox.AddCSSClass("table-content")

	// Add initial rows
	tv.UpdateRows(windows)

	// Create scrolled window
	scrolledWindow := gtk.NewScrolledWindow()
	scrolledWindow.SetChild(tv.ContentBox)
	scrolledWindow.SetVExpand(true)
	scrolledWindow.SetHExpand(true)

	tv.Container.Append(scrolledWindow)

	return tv
}

// UpdateRows updates the table with new window list
func (tv *TableView) UpdateRows(windows []tmux.TmuxWindow) {
	// Remove all existing rows
	for child := tv.ContentBox.FirstChild(); child != nil; child = tv.ContentBox.FirstChild() {
		tv.ContentBox.Remove(child)
	}

	// Store displayed windows for hotkey access
	tv.displayedWindows = windows

	// Add new rows
	for i, window := range windows {
		rowBox := gtk.NewBox(gtk.OrientationHorizontal, 0)
		rowBox.AddCSSClass("table-row")
		if i%2 == 0 {
			rowBox.AddCSSClass("table-row-even")
		} else {
			rowBox.AddCSSClass("table-row-odd")
		}
		rowBox.SetHomogeneous(false)

		// ID label (1-9, then blank for 10+)
		idText := ""
		if i < 9 {
			idText = string(rune('1' + i))
		}
		idLabel := gtk.NewLabel(idText)
		idLabel.SetXAlign(0.5)
		idLabel.SetSizeRequest(40, -1)
		idLabel.AddCSSClass("table-cell")
		idLabel.AddCSSClass("id-cell")

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

		activityLabel := gtk.NewLabel(tmux.FormatTimestamp(window.Activity))
		activityLabel.SetXAlign(0)
		activityLabel.SetSizeRequest(200, -1)
		activityLabel.AddCSSClass("table-cell")

		rowBox.Append(idLabel)
		rowBox.Append(sessionLabel)
		rowBox.Append(indexLabel)
		rowBox.Append(nameLabel)
		rowBox.Append(activityLabel)

		tv.ContentBox.Append(rowBox)
	}
}

// GetWindowByID returns the window at the given index (1-based)
func (tv *TableView) GetWindowByID(id int) *tmux.TmuxWindow {
	if id < 1 || id > len(tv.displayedWindows) {
		return nil
	}
	return &tv.displayedWindows[id-1]
}
