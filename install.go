package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadFile(file string) (string, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/overdodactyl/ShadowFox/master/" + file)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func backUp(path string) error {
	exists, err := pathExists(path)
	if err != nil {
		return err
	}
	if exists {
		err := os.Rename(path, path+".old")
		if err != nil {
			return err
		}
	}
	return nil
}

func createFile(path string) error {
	exists, err := pathExists(path)
	if err != nil {
		return err
	}
	if !exists {
		err := ioutil.WriteFile(path, nil, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func createDir(path string) error {
	pathExists, err := pathExists(path)
	if err != nil {
		return err
	}
	if !pathExists {
		err := os.Mkdir(path, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func addColorOverrides(source, colors string) string {
	startI := strings.Index(source, "--start-indicator-for-updater-scripts: black;")
	endI := strings.Index(source, "--end-indicator-for-updater-scripts: black;") + 43
	return source[:startI] + colors + source[endI:]
}

func install(profilePath string) (string, error) {
	// Helper variables to keep things DRY
	chromePath := filepath.Join(profilePath, "chrome")
	customPath := filepath.Join(chromePath, "ShadowFox_customization")
	userChromePath := filepath.Join(chromePath, "userChrome.css")
	userContentPath := filepath.Join(chromePath, "userContent.css")

	// Create dirs
	if err := createDir(chromePath); err != nil {
		return "Couln't create chrome folder", err
	}
	if err := createDir(customPath); err != nil {
		return "Couln't create ShadowFox_customization folder", err
	}

	// Create customization files
	if err := createFile(filepath.Join(customPath, "colorOverrides.css")); err != nil {
		return "Couln't create colorOverrides.css", err
	}
	if err := createFile(filepath.Join(customPath, "internal_UUIDs.txt")); err != nil {
		return "Couln't create internal_UUIDs.txt", err
	}
	if err := createFile(filepath.Join(customPath, "userContent_customization.css")); err != nil {
		return "Couln't create userContent_customization.css", err
	}
	if err := createFile(filepath.Join(customPath, "userChrome_customization.css")); err != nil {
		return "Couln't create userChrome_customization.css", err
	}

	// Download files
	userChrome, err := downloadFile("userChrome.css")
	if err != nil {
		return "userChrome.css couln't be downloaded", err
	}
	userContent, err := downloadFile("userContent.css")
	if err != nil {
		return "userContent.css couln't be downloaded", err
	}

	// Backup old files
	if err := backUp(userChromePath); err != nil {
		return "Couln't backup userChrome.css", err
	}
	if err := backUp(userContentPath); err != nil {
		return "Couln't backup userContent.css", err
	}

	colors, err := ioutil.ReadFile(filepath.Join(customPath, "colorOverrides.css"))
	if err != nil {
		return "Couln't read colorOverrides.css", err
	}
	if len(colors) != 0 {
		userChrome = addColorOverrides(userChrome, string(colors))
		userContent = addColorOverrides(userContent, string(colors))
	}

	chromeCustom, err := ioutil.ReadFile(filepath.Join(customPath, "userChrome_customization.css"))
	if err != nil {
		return "Couln't read userChrome_customization.css", err
	}
	if len(chromeCustom) != 0 {
		userChrome = userChrome + string(chromeCustom)
	}

	contentCustom, err := ioutil.ReadFile(filepath.Join(customPath, "userContent_customization.css"))
	if err != nil {
		return "Couln't read userContent_customization.css", err
	}
	if len(chromeCustom) != 0 {
		userContent = userContent + string(contentCustom)
	}

	// Write new files
	if err := ioutil.WriteFile(userChromePath, []byte(userChrome), 0644); err != nil {
		return "Couln't write userChrome.css to file", err
	}
	if err := ioutil.WriteFile(userContentPath, []byte(userContent), 0644); err != nil {
		return "Couln't write userContent.css to file", err
	}

	return "", nil
}
