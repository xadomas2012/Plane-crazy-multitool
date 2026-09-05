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

	expectedPrefix := ""

	switch runtime.GOOS {
	case "linux":
		if runtime.GOARCH == "amd64" {
			expectedPrefix = "PC-Multitool-Linux-x64-"
		}

	case "windows":
		if runtime.GOARCH == "amd64" {
			expectedPrefix = "PC-Multitool-Windows-x64-"
		}

	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			expectedPrefix = "PC-Multitool-macOS-x64-"
		case "arm64":
			expectedPrefix = "PC-Multitool-macOS-arm64-"
		}
	}

	if expectedPrefix == "" {
		return githubAsset{},
			fmt.Errorf(
				"unsupported platform: %s/%s",
				runtime.GOOS,
				runtime.GOARCH,
			)
	}

	for _, asset := range assets {
		if strings.HasPrefix(asset.Name, expectedPrefix) &&
			strings.HasSuffix(asset.Name, ".zip") {
			return asset, nil
		}
	}

	return githubAsset{},
		fmt.Errorf(
			"release does not contain platform ZIP starting with %q",
			expectedPrefix,
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

	targetPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	targetPath, err = filepath.Abs(targetPath)
	if err != nil {
		return "", err
	}

	updateDir :=
		filepath.Join(
			filepath.Dir(targetPath),
			".update",
		)

	if err := os.MkdirAll(updateDir, 0755); err != nil {
		return "", fmt.Errorf(
			"cannot create update directory: %w",
			err,
		)
	}

	stagedPath :=
		filepath.Join(
			updateDir,
			"new-"+asset.Name,
		)

	tempPath :=
		filepath.Join(
			updateDir,
			"."+asset.Name+".tmp",
		)

	_ = os.Remove(stagedPath)
	_ = os.Remove(tempPath)

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

	if runtime.GOOS != "windows" {
		if err := os.Chmod(tempPath, 0755); err != nil {
			os.Remove(tempPath)
			return "", err
		}
	}

	if err := os.Rename(tempPath, stagedPath); err != nil {
		os.Remove(tempPath)
		return "",
			fmt.Errorf(
				"cannot stage update: %w",
				err,
			)
	}

	return stagedPath, nil
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

	targetPath, err := os.Executable()
	if err != nil {
		return err
	}

	targetPath, err = filepath.Abs(targetPath)
	if err != nil {
		return err
	}

	updaterName :=
		"PC-Gear-Calculator-Updater"

	if runtime.GOOS == "windows" {
		updaterName += ".exe"
	}

	updaterPath :=
		filepath.Join(
			filepath.Dir(targetPath),
			".update",
			updaterName,
		)

	if _, err := os.Stat(updaterPath); err != nil {
		return fmt.Errorf(
			"bundled updater not found: %w",
			err,
		)
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(updaterPath, 0755); err != nil {
			return err
		}
	}

	pid :=
		strconv.Itoa(
			os.Getpid(),
		)

	cmd :=
		exec.Command(
			updaterPath,
			downloadedPath,
			targetPath,
			pid,
		)

	cmd.Dir =
		filepath.Dir(
			targetPath,
		)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf(
			"unable to start updater: %w",
			err,
		)
	}

	return nil
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
