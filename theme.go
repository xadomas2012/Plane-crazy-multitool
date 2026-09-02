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
	"Catppuccin Latte",
	"Catppuccin Frappé",
	"Catppuccin Macchiato",
	"Catppuccin Mocha",
	"Nord",
	"Gruvbox",
	"Solid",
}

var themeKeys = []string{
	"catppuccin-latte",
	"catppuccin-frappe",
	"catppuccin-macchiato",
	"catppuccin-mocha",
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
	case "catppuccin-latte":
		return palette{
			background: lipgloss.Color("#eff1f5"),
			surface:    lipgloss.Color("#e6e9ef"),
			panel:      lipgloss.Color("#ccd0da"),
			text:       lipgloss.Color("#4c4f69"),
			muted:      lipgloss.Color("#6c6f85"),
			border:     lipgloss.Color("#9ca0b0"),
			accent:     accent,
		}

	case "catppuccin-frappe":
		return palette{
			background: lipgloss.Color("#303446"),
			surface:    lipgloss.Color("#292c3c"),
			panel:      lipgloss.Color("#414559"),
			text:       lipgloss.Color("#c6d0f5"),
			muted:      lipgloss.Color("#a5adce"),
			border:     lipgloss.Color("#737994"),
			accent:     accent,
		}

	case "catppuccin-macchiato":
		return palette{
			background: lipgloss.Color("#24273a"),
			surface:    lipgloss.Color("#1e2030"),
			panel:      lipgloss.Color("#363a4f"),
			text:       lipgloss.Color("#cad3f5"),
			muted:      lipgloss.Color("#a5adcb"),
			border:     lipgloss.Color("#6e738d"),
			accent:     accent,
		}

	case "catppuccin-mocha":
		return palette{
			background: lipgloss.Color("#11111b"),
			surface:    lipgloss.Color("#181825"),
			panel:      lipgloss.Color("#1e1e2e"),
			text:       lipgloss.Color("#cdd6f4"),
			muted:      lipgloss.Color("#7f849c"),
			border:     lipgloss.Color("#45475a"),
			accent:     accent,
		}

	case "catppuccin":
		// Legacy config value: keep old configs working.
		return palette{
			background: lipgloss.Color("#1e1e2e"),
			surface:    lipgloss.Color("#181825"),
			panel:      lipgloss.Color("#313244"),
			text:       lipgloss.Color("#cdd6f4"),
			muted:      lipgloss.Color("#a6adc8"),
			border:     lipgloss.Color("#585b70"),
			accent:     accent,
		}

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
