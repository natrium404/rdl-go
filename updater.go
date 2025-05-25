package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
)

type GitHubRelease struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type ReleseBinary struct {
	Version string
	Asset   Asset
}

type Asset struct {
	DownloadURL string `json:"browser_download_url"`
	Name        string `json:"name"`
}

type VersionInfo struct {
	CurrentVersion string    `json:"current_version"`
	LatestVersion  string    `json:"latest_version"`
	LastChecked    time.Time `json:"last_checked"`
	Asset          Asset     `json:"asset"`
}

func cachedVersion() (*VersionInfo, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	programCacheDir := path.Join(userCacheDir, "rdl")
	err = os.MkdirAll(programCacheDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	var version VersionInfo
	cacheFile := path.Join(programCacheDir, "version.json")

	if _, err := os.Stat(cacheFile); err == nil {
		data, err := os.ReadFile(cacheFile)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &version); err != nil {
			return nil, err
		}

		if time.Since(version.LastChecked) >= 12*time.Hour {
			latestVersion, err := latestVersionGithub()
			if err != nil {
				return nil, err
			}
			version.LatestVersion = latestVersion.Version
			version.Asset = latestVersion.Asset
		}
	} else if os.IsNotExist(err) {
		latestVersion, err := latestVersionGithub()
		if err != nil {
			return nil, err
		}

		version = VersionInfo{
			CurrentVersion: VERSION,
			LatestVersion:  latestVersion.Version,
			LastChecked:    time.Now().UTC(),
			Asset:          latestVersion.Asset,
		}

		data, err := json.MarshalIndent(version, "", "  ")
		if err != nil {
			return nil, err
		}

		if err := os.WriteFile(cacheFile, data, 0644); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return &version, nil
}

func latestVersionGithub() (*ReleseBinary, error) {
	repo := "natrium404/rdl-go"
	resp, err := http.Get("https://api.github.com/repos/" + repo + "/releases/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	binaryName := getBinaryName()

	var GithubRelease GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&GithubRelease); err != nil {
		return nil, err
	}

	var binary ReleseBinary

	for _, asset := range GithubRelease.Assets {
		if asset.Name == binaryName {
			binary = ReleseBinary{
				Version: GithubRelease.TagName,
				Asset:   asset,
			}
			break
		}
	}

	return &binary, nil
}

func getBinaryName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	switch arch {
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "arm64"
	}

	switch os {
	case "windows":
		return fmt.Sprintf("rdl_windows_%s.exe", arch)
	case "linux":
		return fmt.Sprintf("rdl_linux_%s", arch)
	default:
		return "unsupported_platform"
	}
}

func downloadAndReplace(release *Asset) error {
	res, err := http.Get(release.DownloadURL)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	tmpPath := exePath + ".tmp"

	output, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer output.Close()

	if _, err := io.Copy(output, res.Body); err != nil {
		return err
	}

	backupPath := exePath + ".bak"
	if err := os.Rename(exePath, backupPath); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, exePath); err != nil {
		return err
	}

	if err := os.Remove(backupPath); err != nil {
		return err
	}

	return nil
}
