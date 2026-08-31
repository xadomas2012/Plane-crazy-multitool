package main

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type palette struct {
	background color.Color
	surface    color.Color
	panel      color.Color
	text       color.Color
	muted      color.Color
	border     color.Color
	accent     color.Color
}

var themeNames = []string{
	"Catppuccin Mocha",
	"Nord",
	"Gruvbox",
	"Solid",
}

var themeKeys = []string{
	"catppuccin",
	"nord",
	"gruvbox",
	"solid",
}

var accentNames = []string{
	"Green",
	"Blue",
	"Purple",
	"Red",
	"Orange",
	"Cyan",
	"Pink",
}

var accentKeys = []string{
	"green",
	"blue",
	"purple",
	"red",
	"orange",
	"cyan",
	"pink",
}

var accentColors = map[string]color.Color{
	"green":  lipgloss.Color("#a6e3a1"),
	"blue":   lipgloss.Color("#89b4fa"),
	"purple": lipgloss.Color("#cba6f7"),
	"red":    lipgloss.Color("#f38ba8"),
	"orange": lipgloss.Color("#fab387"),
	"cyan":   lipgloss.Color("#89dceb"),
	"pink":   lipgloss.Color("#f5c2e7"),
}

func getPalette(themeKey, accentKey string) palette {
	accent, ok := accentColors[accentKey]
	if !ok {
		accent = accentColors["green"]
	}

	switch themeKey {
	case "nord":
		return palette{
			background: lipgloss.Color("#2e3440"),
			surface:    lipgloss.Color("#3b4252"),
			panel:      lipgloss.Color("#434c5e"),
			text:       lipgloss.Color("#eceff4"),
			muted:      lipgloss.Color("#88c0d0"),
			border:     lipgloss.Color("#4c566a"),
			accent:     accent,
		}

	case "gruvbox":
		return palette{
			background: lipgloss.Color("#1d2021"),
			surface:    lipgloss.Color("#282828"),
			panel:      lipgloss.Color("#3c3836"),
			text:       lipgloss.Color("#ebdbb2"),
			muted:      lipgloss.Color("#a89984"),
			border:     lipgloss.Color("#504945"),
			accent:     accent,
		}

	case "solid":
		return palette{
			background: lipgloss.Color("#080808"),
			surface:    lipgloss.Color("#101010"),
			panel:      lipgloss.Color("#151515"),
			text:       lipgloss.Color("#eeeeee"),
			muted:      lipgloss.Color("#777777"),
			border:     lipgloss.Color("#303030"),
			accent:     accent,
		}

	default:
		return palette{
			background: lipgloss.Color("#11111b"),
			surface:    lipgloss.Color("#181825"),
			panel:      lipgloss.Color("#1e1e2e"),
			text:       lipgloss.Color("#cdd6f4"),
			muted:      lipgloss.Color("#7f849c"),
			border:     lipgloss.Color("#45475a"),
			accent:     accent,
		}
	}
}
