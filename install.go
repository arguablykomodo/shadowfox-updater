package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const darkThemeConfig = `
user_pref("lightweightThemes.selectedThemeID", "firefox-compact-dark@mozilla.org");
user_pref("devtools.theme", "dark");`

func uninstall(profile string) (string, error) {
	err := os.RemoveAll(filepath.Join(profile, "chrome", "ShadowFox_customization"))
	if err != nil {
		return "Couldn't delete ShadowFox_customization", err
	}
	err = os.Remove(filepath.Join(profile, "chrome", "userChrome.css"))
	if err != nil {
		return "Couldn't delete userChrome.css", err
	}
	err = os.Remove(filepath.Join(profile, "chrome", "userContent.css"))
	if err != nil {
		return "Couldn't delete userContent.css", err
	}
	return "", nil
}

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

func pathExists(path string) (bool, bool, error) {
	stats, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}
	return true, stats.IsDir(), nil
}

func backUp(path string) error {
	exists, _, err := pathExists(path)
	if err != nil {
		return err
	}
	if exists {
		err := os.Rename(path, path+time.Now().Format(".2006-01-02-15-04-05.backup"))
		if err != nil {
			return err
		}
	}
	return nil
}

func createFile(path string) error {
	exists, _, err := pathExists(path)
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
	pathExists, _, err := pathExists(path)
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

func install(profilePath string, generateUUIDs bool, setTheme bool) (string, error) {
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
	if err := createFile(filepath.Join(customPath, "internal_UUIDs.txt")); err != nil {
		return "Couldn't create internal_UUIDs.txt", err
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
	uuidFile, err := ioutil.ReadFile(filepath.Join(customPath, "internal_UUIDs.txt"))
	if err != nil {
		return "Couldn't read internal_UUIDs.txt", err
	}
	if generateUUIDs {
		err := backUp(filepath.Join(customPath, "internal_UUIDs.txt"))
		if err != nil {
			return "Couldn't backup internal_UUIDs.txt", err
		}
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
		newUUIDFile := ""
		for key, value := range prefsJSON {
			newUUIDFile += key + "=" + value + "\n"
			userContent = strings.Replace(userContent, key, value, -1)
		}
		if err := ioutil.WriteFile(filepath.Join(customPath, "internal_UUIDs.txt"), []byte(newUUIDFile), 0644); err != nil {
			return "Couldn't write internal_UUIDs.txt to file", err
		}
	} else {
		pairs := regexp.MustCompile("(.+)=(.+)").FindAllStringSubmatch(string(uuidFile), -1)
		for _, key := range pairs {
			userContent = strings.Replace(userContent, key[1], key[2], -1)
		}
	}

	// Set dark theme
	if setTheme {
		userJs := filepath.Join(profilePath, "user.js")
		userJsContent := []byte{}

		exists, _, err := pathExists(userJs)
		if exists {
			userJsContent, err = ioutil.ReadFile(userJs)
			if err != nil {
				return "Couldn't read user.js", err
			}
		} else {
			err = createFile(userJs)
			if err != nil {
				return "Couldn't create user.js", err
			}
		}

		if !strings.Contains(string(userJsContent), darkThemeConfig) {
			err = backUp(userJs)
			if err != nil {
				return "Couldn't backup user.js", err
			}

			userJsContent = append(userJsContent, []byte(darkThemeConfig)...)

			if err := ioutil.WriteFile(userJs, userJsContent, 0644); err != nil {
				return "Couldn't write user.js", err
			}
		}
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
