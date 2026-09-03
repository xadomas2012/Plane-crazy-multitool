package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var Version = "dev"

const (
	githubRepo       = "xadomas2012/Plane-crazy-multitool"
	latestReleaseAPI = "https://api.github.com/repos/" +
		githubRepo +
		"/releases/latest"
)

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Name    string        `json:"name"`
	Body    string        `json:"body"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Digest             string `json:"digest"`
	Size               int64  `json:"size"`
}

type updateInfo struct {
	CurrentVersion string
	LatestVersion  string
	ReleaseName    string
	ReleaseNotes   string
	Asset          githubAsset
}

func fetchLatestRelease() (githubRelease, error) {
	client :=
		&http.Client{
			Timeout: 10 * time.Second,
		}

	req, err :=
		http.NewRequest(
			http.MethodGet,
			latestReleaseAPI,
			nil,
		)

	if err != nil {
		return githubRelease{}, err
	}

	req.Header.Set(
		"Accept",
		"application/vnd.github+json",
	)

	req.Header.Set(
		"X-GitHub-Api-Version",
		"2022-11-28",
	)

	resp, err :=
		client.Do(req)

	if err != nil {
		return githubRelease{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {

		return githubRelease{},
			fmt.Errorf(
				"github API returned HTTP %d",
				resp.StatusCode,
			)
	}

	var release githubRelease

	if err :=
		json.NewDecoder(
			resp.Body,
		).Decode(&release); err != nil {

		return githubRelease{}, err
	}

	if strings.TrimSpace(release.TagName) == "" {
		return githubRelease{},
			errors.New(
				"GitHub release has no tag",
			)
	}

	return release, nil
}

func getLatestUpdate() (updateInfo, error) {
	if Version == "dev" {
		return updateInfo{},
			errors.New(
				"running development build",
			)
	}

	release, err :=
		fetchLatestRelease()

	if err != nil {
		return updateInfo{}, err
	}

	if !isNewerVersion(
		release.TagName,
		Version,
	) {
		return updateInfo{
			CurrentVersion: Version,
			LatestVersion:  release.TagName,
			ReleaseName:    release.Name,
			ReleaseNotes:   release.Body,
		}, nil
	}

	asset, err :=
		findPlatformAsset(
			release.Assets,
		)

	if err != nil {
		return updateInfo{}, err
	}

	return updateInfo{
		CurrentVersion: Version,
		LatestVersion:  release.TagName,
		ReleaseName:    release.Name,
		ReleaseNotes:   release.Body,
		Asset:          asset,
	}, nil
}

func findPlatformAsset(
	assets []githubAsset,
) (githubAsset, error) {

	expected := ""

	switch runtime.GOOS {

	case "linux":

		if runtime.GOARCH == "amd64" {
			expected =
				"PC-Gear-Calculator-Linux-x64"
		}

	case "windows":

		if runtime.GOARCH == "amd64" {
			expected =
				"PC-Gear-Calculator-Windows-x64.exe"
		}

	case "darwin":

		switch runtime.GOARCH {

		case "amd64":
			expected =
				"PC-Gear-Calculator-macOS-x64"

		case "arm64":
			expected =
				"PC-Gear-Calculator-macOS-arm64"
		}
	}

	if expected == "" {
		return githubAsset{},
			fmt.Errorf(
				"unsupported platform: %s/%s",
				runtime.GOOS,
				runtime.GOARCH,
			)
	}

	for _, asset := range assets {
		if asset.Name == expected {
			return asset, nil
		}
	}

	return githubAsset{},
		fmt.Errorf(
			"release does not contain asset %q",
			expected,
		)
}

func downloadUpdate(
	asset githubAsset,
) (string, error) {

	if asset.BrowserDownloadURL == "" {
		return "",
			errors.New(
				"release asset has no download URL",
			)
	}

	client :=
		&http.Client{
			Timeout: 5 * time.Minute,
		}

	resp, err :=
		client.Get(
			asset.BrowserDownloadURL,
		)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {

		return "",
			fmt.Errorf(
				"download returned HTTP %d",
				resp.StatusCode,
			)
	}

	tempDir :=
		os.TempDir()

	tempPath :=
		filepath.Join(
			tempDir,
			"pc-multitool-update-"+asset.Name,
		)

	file, err :=
		os.Create(tempPath)

	if err != nil {
		return "", err
	}

	_, copyErr :=
		io.Copy(
			file,
			resp.Body,
		)

	closeErr :=
		file.Close()

	if copyErr != nil {
		os.Remove(tempPath)
		return "", copyErr
	}

	if closeErr != nil {
		os.Remove(tempPath)
		return "", closeErr
	}

	if asset.Size > 0 {

		info, err :=
			os.Stat(tempPath)

		if err != nil {
			os.Remove(tempPath)
			return "", err
		}

		if info.Size() != asset.Size {

			os.Remove(tempPath)

			return "",
				fmt.Errorf(
					"downloaded size mismatch: got %d, expected %d",
					info.Size(),
					asset.Size,
				)
		}
	}

	if err :=
		verifyDownloadedFile(
			tempPath,
			asset.Digest,
		); err != nil {

		os.Remove(tempPath)
		return "", err
	}

	return tempPath, nil
}

func verifyDownloadedFile(
	path string,
	expectedDigest string,
) error {

	expectedDigest =
		strings.TrimSpace(
			expectedDigest,
		)

	expectedDigest =
		strings.TrimPrefix(
			expectedDigest,
			"sha256:",
		)

	if expectedDigest == "" {
		return errors.New(
			"release asset has no SHA-256 digest",
		)
	}

	file, err :=
		os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	hash :=
		sha256.New()

	if _, err :=
		io.Copy(
			hash,
			file,
		); err != nil {

		return err
	}

	actualDigest :=
		hex.EncodeToString(
			hash.Sum(nil),
		)

	if !strings.EqualFold(
		actualDigest,
		expectedDigest,
	) {
		return fmt.Errorf(
			"SHA-256 mismatch: got %s, expected %s",
			actualDigest,
			expectedDigest,
		)
	}

	return nil
}

func startUpdateHelper(
	downloadedPath string,
) error {

	targetPath, err :=
		os.Executable()

	if err != nil {
		return err
	}

	targetPath, err =
		filepath.Abs(
			targetPath,
		)

	if err != nil {
		return err
	}

	pid :=
		strconv.Itoa(
			os.Getpid(),
		)

	switch runtime.GOOS {

	case "windows":

		helperPath :=
			filepath.Join(
				os.TempDir(),
				fmt.Sprintf(
					"pc-multitool-update-%d.cmd",
					time.Now().UnixNano(),
				),
			)

		script :=
			`@echo off
set "NEW=%~1"
set "TARGET=%~2"
set "PID=%~3"

:wait
tasklist /FI "PID eq %PID%" 2>NUL | find "%PID%" >NUL
if not errorlevel 1 (
    timeout /t 1 /nobreak >NUL
    goto wait
)

:replace
copy /Y "%NEW%" "%TARGET%" >NUL 2>&1
if errorlevel 1 (
    timeout /t 1 /nobreak >NUL
    goto replace
)

del "%NEW%" >NUL 2>&1
start "" "%TARGET%"
del "%~f0" >NUL 2>&1
`

		if err :=
			os.WriteFile(
				helperPath,
				[]byte(script),
				0700,
			); err != nil {
			return err
		}

		cmd :=
			exec.Command(
				"cmd.exe",
				"/C",
				"start",
				"",
				"/B",
				helperPath,
				downloadedPath,
				targetPath,
				pid,
			)

		if err :=
			cmd.Start(); err != nil {
			os.Remove(helperPath)
			return err
		}

		return nil

	case "linux", "darwin":

		helperPath :=
			filepath.Join(
				os.TempDir(),
				fmt.Sprintf(
					"pc-multitool-update-%d.sh",
					time.Now().UnixNano(),
				),
			)

		script :=
			`#!/bin/sh

NEW="$1"
TARGET="$2"

sleep 0.5

while ! mv -f "$NEW" "$TARGET" 2>/dev/null
do
    sleep 0.2
done

chmod +x "$TARGET"

"$TARGET" &

rm -f "$0"
`

		if err :=
			os.WriteFile(
				helperPath,
				[]byte(script),
				0700,
			); err != nil {
			return err
		}

		cmd :=
			exec.Command(
				"sh",
				helperPath,
				downloadedPath,
				targetPath,
			)

		if err :=
			cmd.Start(); err != nil {
			os.Remove(helperPath)
			return err
		}

		return nil

	default:

		return fmt.Errorf(
			"unsupported update platform: %s",
			runtime.GOOS,
		)
	}
}

func checkForUpdate() {
	if Version == "dev" {
		return
	}

	info, err :=
		getLatestUpdate()

	if err != nil {
		return
	}

	if !isNewerVersion(
		info.LatestVersion,
		Version,
	) {
		return
	}

	// Keep the existing startup notification behavior
	// for now. The TUI will use getLatestUpdate()
	// directly for the actual updater flow.
	_ = os.WriteFile(
		filepath.Join(
			os.TempDir(),
			"pc-multitool-update-available",
		),
		[]byte(
			info.LatestVersion,
		),
		0600,
	)
}

type parsedVersion struct {
	major      int
	minor      int
	patch      int
	preRelease bool
	preID      string
	preNumber  int
}

func isNewerVersion(latest, current string) bool {
	latestVersion, latestOK := parseVersion(latest)
	currentVersion, currentOK := parseVersion(current)

	if !latestOK || !currentOK {
		return false
	}

	if latestVersion.major != currentVersion.major {
		return latestVersion.major > currentVersion.major
	}

	if latestVersion.minor != currentVersion.minor {
		return latestVersion.minor > currentVersion.minor
	}

	if latestVersion.patch != currentVersion.patch {
		return latestVersion.patch > currentVersion.patch
	}

	if latestVersion.preRelease != currentVersion.preRelease {
		return !latestVersion.preRelease
	}

	if !latestVersion.preRelease {
		return false
	}

	preRank := func(id string) int {
		switch strings.ToLower(id) {
		case "alpha":
			return 1
		case "beta":
			return 2
		case "rc":
			return 3
		default:
			return 0
		}
	}

	latestRank := preRank(latestVersion.preID)
	currentRank := preRank(currentVersion.preID)

	if latestRank != currentRank {
		return latestRank > currentRank
	}

	return latestVersion.preNumber > currentVersion.preNumber
}

func parseVersion(version string) (parsedVersion, bool) {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")

	parts := strings.SplitN(version, "-", 2)
	core := strings.Split(parts[0], ".")

	if len(core) != 3 {
		return parsedVersion{}, false
	}

	major, err := strconv.Atoi(core[0])
	if err != nil {
		return parsedVersion{}, false
	}

	minor, err := strconv.Atoi(core[1])
	if err != nil {
		return parsedVersion{}, false
	}

	patch, err := strconv.Atoi(core[2])
	if err != nil {
		return parsedVersion{}, false
	}

	result := parsedVersion{
		major: major,
		minor: minor,
		patch: patch,
	}

	if len(parts) == 1 {
		return result, true
	}

	pre := strings.Split(parts[1], ".")

	if len(pre) != 2 || pre[0] == "" {
		return parsedVersion{}, false
	}

	preNumber, err := strconv.Atoi(pre[1])
	if err != nil || preNumber < 0 {
		return parsedVersion{}, false
	}

	result.preRelease = true
	result.preID = pre[0]
	result.preNumber = preNumber

	return result, true
}
