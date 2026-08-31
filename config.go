package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	Theme  string `json:"theme"`
	Accent string `json:"accent"`
}

func configPath() string {
	dir, err := os.UserConfigDir()
	if err == nil {
		return filepath.Join(dir, "PC-Gear-Calculator", "config.json")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "pc-gear-calculator.json"
	}

	return filepath.Join(home, ".pc-gear-calculator.json")
}

func loadConfig() config {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return config{
			Theme:  "catppuccin",
			Accent: "green",
		}
	}

	var cfg config

	if err := json.Unmarshal(data, &cfg); err != nil {
		return config{
			Theme:  "catppuccin",
			Accent: "green",
		}
	}

	if cfg.Theme == "" {
		cfg.Theme = "catppuccin"
	}

	if cfg.Accent == "" {
		cfg.Accent = "green"
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
