package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var Version = "dev"

const latestReleaseURL = "https://github.com/xadomas2012/Plane-crazy-multitool/releases/latest"

func checkForUpdate() {
	if Version == "dev" {
		return
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(latestReleaseURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return
	}

	latest := strings.TrimPrefix(
		resp.Request.URL.Path,
		"/xadomas2012/Plane-crazy-multitool/releases/tag/",
	)

	latest = strings.TrimSpace(latest)

	if latest == "" || !isNewerVersion(latest, Version) {
		return
	}

	_ = exec.Command(
		"notify-send",
		"-u", "normal",
		"-a", "PC Gear Calculator",
		"New version available",
		fmt.Sprintf(
			"PC Gear Calculator %s is available.\nYou are running %s.\nPlease update from GitHub.",
			latest,
			Version,
		),
	).Run()
}

func isNewerVersion(latest, current string) bool {
	latestParts := parseVersion(latest)
	currentParts := parseVersion(current)

	for i := 0; i < 3; i++ {
		if latestParts[i] > currentParts[i] {
			return true
		}

		if latestParts[i] < currentParts[i] {
			return false
		}
	}

	return false
}

func parseVersion(version string) [3]int {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")

	parts := strings.Split(version, ".")
	var result [3]int

	for i := 0; i < len(parts) && i < 3; i++ {
		n, err := strconv.Atoi(parts[i])
		if err != nil {
			return [3]int{}
		}

		result[i] = n
	}

	return result
}
