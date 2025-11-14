package tmux

import (
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TmuxWindow struct {
	Session  string
	Index    string
	Name     string
	Activity string
}

// GetWindows retrieves all tmux windows sorted by most recent activity
func GetWindows() []TmuxWindow {
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

// FormatTimestamp converts Unix timestamp to human-readable format
func FormatTimestamp(ts string) string {
	timestamp, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return ts
	}

	t := time.Unix(timestamp, 0)

	// Format as "15:04 2 Jan"
	return t.Format("15:04 2 Jan")
}

// FormatForFzf formats a window for fzf input
func (w TmuxWindow) FormatForFzf() string {
	return w.Session + ":" + w.Index + "  " + w.Name + "  " + FormatTimestamp(w.Activity)
}

// SwitchToWindow switches to the specified tmux window
func SwitchToWindow(session, index string) error {
	target := session + ":" + index
	cmd := exec.Command("tmux", "switch-client", "-t", target)
	return cmd.Run()
}
