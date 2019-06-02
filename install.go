package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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

func backUp(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	err = os.Rename(path, path+time.Now().Format(".2006-01-02-15-04-05.backup"))
	if err != nil {
		return err
	}
	return nil
}

func readFile(path string) (string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := ioutil.WriteFile(path, nil, 0644)
		if err != nil {
			return "", err
		}
		return "", nil
	}
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func uninstall(profile string) error {
	for _, file := range []string{
		"userChrome.css",
		"userContent.css",
	} {
		path := filepath.Join(profile, "chrome", file)
		if err := backUp(path); err != nil {
			return fmt.Errorf("Couldn't backup %s: %s", file, err)
		}
	}
	return nil
}

func install(profilePath string, generateUUIDs bool, setTheme bool) error {
	chromePath := filepath.Join(profilePath, "chrome")
	customPath := filepath.Join(chromePath, "ShadowFox_customization")

	if err := os.MkdirAll(customPath, 0700); err != nil {
		return fmt.Errorf("Couldn't create folders: %s", err)
	}

	colors, err := readFile(filepath.Join(customPath, "colorOverrides.css"))
	if err != nil {
		return fmt.Errorf("Couldn't read colorOverrides.css: %s", err)
	}

	if generateUUIDs {
		err := backUp(filepath.Join(customPath, "internal_UUIDs.txt"))
		if err != nil {
			return fmt.Errorf("Couldn't backup internal_UUIDs.txt: %s", err)
		}

		prefs, err := readFile(filepath.Join(profilePath, "prefs.js"))
		if err != nil {
			return fmt.Errorf("Couldn't read prefs.js: %s", err)
		}

		regex := regexp.MustCompile(`\\\"(.+?)\\\":\\\"(.{8}-.{4}-.{4}-.{4}-.{12})\\\"`)
		matches := regex.FindAllStringSubmatch(prefs, -1)
		output := ""
		for _, match := range matches {
			output += match[1] + "=" + match[2] + "\n"
		}

		if err := ioutil.WriteFile(filepath.Join(customPath, "internal_UUIDs.txt"), []byte(output), 0644); err != nil {
			return fmt.Errorf("Couldn't write internal_UUIDs.txt: %s", err)
		}
	}

	uuidBytes, err := readFile(filepath.Join(customPath, "internal_UUIDs.txt"))
	if err != nil {
		return fmt.Errorf("Couldn't read internal_UUIDs.txt: %s", err)
	}
	uuids := string(uuidBytes)
	pairs := regexp.MustCompile("(.+)=(.+)").FindAllStringSubmatch(uuids, -1)

	for _, file := range []string{
		"userChrome",
		"userContent",
	} {
		path := filepath.Join(chromePath, file)

		if err := backUp(path + ".css"); err != nil {
			return fmt.Errorf("Couldn't backup %s: %s", file, err)
		}

		contents, err := downloadFile(file + ".css")
		if err != nil {
			return fmt.Errorf("Couldn't download %s: %s", file, err)
		}

		// Add color overrides
		startI := strings.Index(contents, "--start-indicator-for-updater-scripts: black;")
		endI := strings.Index(contents, "--end-indicator-for-updater-scripts: black;") + 43
		contents = contents[:startI] + colors + contents[endI:]

		// Add customizations
		custom, err := readFile(filepath.Join(customPath, file+"_customization.css"))
		if err != nil {
			return fmt.Errorf("Couldn't read %s_customization.css: %s", file, err)
		}
		contents = contents + string(custom)

		// Add UUIDs
		for _, key := range pairs {
			contents = strings.Replace(contents, key[1], key[2], -1)
		}

		// Write file
		if err := ioutil.WriteFile(path+".css", []byte(contents), 0644); err != nil {
			return fmt.Errorf("Couldn't write %s: %s", file, err)
		}
	}

	// Set dark theme
	if setTheme {
		path := filepath.Join(profilePath, "prefs.js")
		prefsContent, err := readFile(path)
		if err != nil {
			return fmt.Errorf("Couldn't read prefs.js: %s", err)
		}

		for key, value := range map[string]string{
			"lightweightThemes.selectedThemeID": "\"firefox-compact-dark@mozilla.org\"",
			"browser.uidensity":                 "1",
			"devtools.theme":                    "\"dark\"",
		} {
			regex := regexp.MustCompile("user_pref(\"" + key + "\", .+);")
			replace := "user_pref(\"" + key + "\", " + value + ");"
			if regex.MatchString(prefsContent) {
				prefsContent = regex.ReplaceAllString(prefsContent, replace)
			} else {
				prefsContent += replace + "\n"
			}
		}

		if err := ioutil.WriteFile(path, []byte(prefsContent), 0644); err != nil {
			return fmt.Errorf("Couldn't write prefs.js: %s", err)
		}
	}

	return nil
}
