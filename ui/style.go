package ui

import (
	_ "embed"
	"log"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed style.css
var styleCSS string

// LoadCSS loads the embedded CSS and applies it globally
func LoadCSS() {
	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		loadCSSProvider(styleCSS),
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}

func loadCSSProvider(content string) *gtk.CSSProvider {
	prov := gtk.NewCSSProvider()
	prov.ConnectParsingError(func(sec *gtk.CSSSection, err error) {
		loc := sec.StartLocation()
		lines := strings.Split(content, "\n")
		log.Printf("CSS error (%v) at line: %q", err, lines[loc.Lines()])
	})
	prov.LoadFromData(content)
	return prov
}
