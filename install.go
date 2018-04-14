package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func downloadFile(file string) ([]byte, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/overdodactyl/ShadowFox/master/" + file)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
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

func install(profilePath string) (string, error) {
	userChrome, err := downloadFile("userChrome.css")
	if err != nil {
		return "userChrome.css couln't be downloaded", err
	}

	userContent, err := downloadFile("userContent.css")
	if err != nil {
		return "userContent.css couln't be downloaded", err
	}

	chromeExists, err := pathExists(filepath.Join(profilePath, "chrome"))
	if err != nil {
		return "", err
	}
	if !chromeExists {
		err := os.Mkdir(filepath.Join(profilePath, "chrome"), 0644)
		if err != nil {
			return "", err
		}
	}

	userChromePath := filepath.Join(profilePath, "chrome", "userChrome.css")
	userContentPath := filepath.Join(profilePath, "chrome", "userContent.css")

	userChromeExits, err := pathExists(userChromePath)
	if err != nil {
		return "", err
	}
	if userChromeExits {
		err := os.Rename(userChromePath, userChromePath+".old")
		if err != nil {
			return "Couln't backup old userChrome.css", err
		}
	}

	userContentExits, err := pathExists(userContentPath)
	if err != nil {
		return "", err
	}
	if userContentExits {
		err := os.Rename(userContentPath, userContentPath+".old")
		if err != nil {
			return "Couln't backup old userContent.css", err
		}
	}

	err = ioutil.WriteFile(userChromePath, userChrome, 0644)
	if err != nil {
		return "Couln't write userChrome.css to file", err
	}

	err = ioutil.WriteFile(userContentPath, userContent, 0644)
	if err != nil {
		return "Couln't write userContent.css to file", err
	}

	return "", nil
}
