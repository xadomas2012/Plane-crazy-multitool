package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	Appearance appearanceConfig `json:"appearance"`
	Layout     layoutConfig     `json:"layout"`
	Reference  referenceConfig  `json:"reference"`
	Calculator calculatorConfig `json:"calculator"`
	Density    string           `json:"density"`
}

type appearanceConfig struct {
	Theme                 string `json:"theme"`
	Accent                string `json:"accent"`
	Background            string `json:"background"`
	TopBarTransparent     bool   `json:"top_bar_transparent"`
	ReferenceTransparent  bool   `json:"reference_transparent"`
	CalculatorTransparent bool   `json:"calculator_transparent"`
	BottomBarTransparent  bool   `json:"bottom_bar_transparent"`
}

type layoutConfig struct {
	Mode           string `json:"mode"`
	ReferenceWidth string `json:"reference_width"`
	Order          string `json:"order"`
}

type referenceConfig struct {
	Enabled         bool `json:"enabled"`
	Teeth           bool `json:"teeth"`
	FullAngle       bool `json:"full_angle"`
	HalfAngle       bool `json:"half_angle"`
	Offset          bool `json:"offset"`
	CompressorValue bool `json:"compressor_value"`
	Compressors     bool `json:"compressors"`
}

type calculatorConfig struct {
	Enabled         bool `json:"enabled"`
	Teeth           bool `json:"teeth"`
	Compressors     bool `json:"compressors"`
	FullAngle       bool `json:"full_angle"`
	HalfAngle       bool `json:"half_angle"`
	Offset          bool `json:"offset"`
	CompressorValue bool `json:"compressor_value"`
	Warnings        bool `json:"warnings"`
}

func defaultConfig() config {
	return config{
		Appearance: appearanceConfig{
			Theme:                 "catppuccin",
			Accent:                "green",
			Background:            "theme",
			TopBarTransparent:     false,
			ReferenceTransparent:  false,
			CalculatorTransparent: false,
			BottomBarTransparent:  false,
		},

		Layout: layoutConfig{
			Mode:           "automatic",
			ReferenceWidth: "balanced",
			Order:          "reference-first",
		},

		Reference: referenceConfig{
			Enabled:         true,
			Teeth:           true,
			FullAngle:       true,
			HalfAngle:       true,
			Offset:          true,
			CompressorValue: true,
			Compressors:     true,
		},

		Calculator: calculatorConfig{
			Enabled:         true,
			Teeth:           true,
			Compressors:     true,
			FullAngle:       true,
			HalfAngle:       true,
			Offset:          true,
			CompressorValue: true,
			Warnings:        true,
		},

		Density: "normal",
	}
}

func configPath() string {
	dir, err := os.UserConfigDir()
	if err == nil {
		return filepath.Join(
			dir,
			"PC-Gear-Calculator",
			"config.json",
		)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "pc-gear-calculator.json"
	}

	return filepath.Join(
		home,
		".pc-gear-calculator.json",
	)
}

func loadConfig() config {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return defaultConfig()
	}

	var raw map[string]json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return defaultConfig()
	}

	// Start with the real defaults, then overlay only settings
	// that actually exist in the saved configuration.
	cfg := defaultConfig()

	if value, ok := raw["appearance"]; ok {
		_ = json.Unmarshal(value, &cfg.Appearance)
	}

	if value, ok := raw["layout"]; ok {
		_ = json.Unmarshal(value, &cfg.Layout)
	}

	if value, ok := raw["reference"]; ok {
		_ = json.Unmarshal(value, &cfg.Reference)
	}

	if value, ok := raw["calculator"]; ok {
		_ = json.Unmarshal(value, &cfg.Calculator)
	}

	if value, ok := raw["density"]; ok {
		_ = json.Unmarshal(value, &cfg.Density)
	}

	return normalizeConfig(cfg)
}

func normalizeConfig(cfg config) config {
	defaults := defaultConfig()

	if cfg.Appearance.Theme == "" {
		cfg.Appearance.Theme = defaults.Appearance.Theme
	}

	if cfg.Appearance.Accent == "" {
		cfg.Appearance.Accent = defaults.Appearance.Accent
	}

	if cfg.Appearance.Background == "" {
		cfg.Appearance.Background = defaults.Appearance.Background
	}

	if cfg.Layout.Mode == "" {
		cfg.Layout.Mode = defaults.Layout.Mode
	}

	if cfg.Layout.ReferenceWidth == "" {
		cfg.Layout.ReferenceWidth = defaults.Layout.ReferenceWidth
	}

	if cfg.Layout.Order == "" {
		cfg.Layout.Order = defaults.Layout.Order
	}

	if cfg.Density == "" {
		cfg.Density = defaults.Density
	}

	return cfg
}

func saveConfig(cfg config) {
	path := configPath()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}

	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return
	}

	_ = os.WriteFile(path, data, 0o644)
}
