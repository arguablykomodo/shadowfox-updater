package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
		err := os.Mkdir(path, 0700)
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
		return "Couldn't create chrome folder", err
	}
	if err := createDir(customPath); err != nil {
		return "Couldn't create ShadowFox_customization folder", err
	}

	// Create customization files
	if err := createFile(filepath.Join(customPath, "colorOverrides.css")); err != nil {
		return "Couldn't create colorOverrides.css", err
	}
	if err := createFile(filepath.Join(customPath, "userContent_customization.css")); err != nil {
		return "Couldn't create userContent_customization.css", err
	}
	if err := createFile(filepath.Join(customPath, "userChrome_customization.css")); err != nil {
		return "Couldn't create userChrome_customization.css", err
	}

	// Download files
	userChrome, err := downloadFile("userChrome.css")
	if err != nil {
		return "userChrome.css Couldn't be downloaded", err
	}
	userContent, err := downloadFile("userContent.css")
	if err != nil {
		return "userContent.css Couldn't be downloaded", err
	}

	// Backup old files
	if err := backUp(userChromePath); err != nil {
		return "Couldn't backup userChrome.css", err
	}
	if err := backUp(userContentPath); err != nil {
		return "Couldn't backup userContent.css", err
	}

	// Add color overrides
	colors, err := ioutil.ReadFile(filepath.Join(customPath, "colorOverrides.css"))
	if err != nil {
		return "Couldn't read colorOverrides.css", err
	}
	if len(colors) != 0 {
		userChrome = addColorOverrides(userChrome, string(colors))
		userContent = addColorOverrides(userContent, string(colors))
	}

	// Add customization files
	chromeCustom, err := ioutil.ReadFile(filepath.Join(customPath, "userChrome_customization.css"))
	if err != nil {
		return "Couldn't read userChrome_customization.css", err
	}
	if len(chromeCustom) != 0 {
		userChrome = userChrome + string(chromeCustom)
	}

	contentCustom, err := ioutil.ReadFile(filepath.Join(customPath, "userContent_customization.css"))
	if err != nil {
		return "Couldn't read userContent_customization.css", err
	}
	if len(chromeCustom) != 0 {
		userContent = userContent + string(contentCustom)
	}

	// Add UUIDs
	prefsFile, err := ioutil.ReadFile(filepath.Join(profilePath, "prefs.js"))
	if err != nil {
		return "Couldn't read prefs.js", err
	}
	prefsString := strings.Replace(
		regexp.MustCompile("extensions\\.webextensions\\.uuids\\\", \\\"(.+)\\\"\\);").
			FindStringSubmatch(string(prefsFile))[1],
		"\\", "", -1,
	)
	var prefsJSON map[string]string
	err = json.Unmarshal([]byte(prefsString), &prefsJSON)
	if err != nil {
		return "Couldn't parse prefs.js", err
	}
	for key, value := range prefsJSON {
		userContent = strings.Replace(userContent, key, value, -1)
	}

	// Write new files
	if err := ioutil.WriteFile(userChromePath, []byte(userChrome), 0644); err != nil {
		return "Couldn't write userChrome.css to file", err
	}
	if err := ioutil.WriteFile(userContentPath, []byte(userContent), 0644); err != nil {
		return "Couldn't write userContent.css to file", err
	}

	return "", nil
}
