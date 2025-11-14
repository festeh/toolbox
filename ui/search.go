package ui

import (
	"bytes"
	"os/exec"
	"strings"
	"toolbox/tmux"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type SearchInput struct {
	Entry      *gtk.Entry
	tableView  *TableView
	allWindows []tmux.TmuxWindow
}

// CreateSearchInput creates a search input field with fzf filtering
func CreateSearchInput(tableView *TableView, windows []tmux.TmuxWindow) *SearchInput {
	si := &SearchInput{
		Entry:      gtk.NewEntry(),
		tableView:  tableView,
		allWindows: windows,
	}

	si.Entry.SetPlaceholderText("Search windows...")
	si.Entry.AddCSSClass("search-input")

	// Connect text change handler
	si.Entry.ConnectChanged(func() {
		query := si.Entry.Text()
		filtered := si.filterWithFzf(query)
		si.tableView.UpdateRows(filtered)
	})

	return si
}

// filterWithFzf uses fzf --filter to fuzzy match windows
func (si *SearchInput) filterWithFzf(query string) []tmux.TmuxWindow {
	// If query is empty, return all windows
	if query == "" {
		return si.allWindows
	}

	// Build input for fzf
	var input strings.Builder
	windowMap := make(map[string]tmux.TmuxWindow)

	for _, window := range si.allWindows {
		line := window.FormatForFzf()
		input.WriteString(line + "\n")
		windowMap[line] = window
	}

	// Run fzf with --filter flag
	cmd := exec.Command("fzf", "--filter="+query)
	cmd.Stdin = strings.NewReader(input.String())

	var out bytes.Buffer
	cmd.Stdout = &out

	// Run fzf (ignore error as it may return non-zero on no matches)
	cmd.Run()

	// Parse output back to windows
	filtered := []tmux.TmuxWindow{}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		if window, ok := windowMap[line]; ok {
			filtered = append(filtered, window)
		}
	}

	return filtered
}
