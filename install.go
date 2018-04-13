package main

import (
	"errors"
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

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	checkErr(err)
	return true
}

func install() error {
	userChrome, err := downloadFile("userChrome.css")
	if err != nil {
		return errors.New("userChrome.css couln't be downloaded: " + err.Error())
	}

	userContent, err := downloadFile("userContent.css")
	if err != nil {
		return errors.New("userContent.css couln't be downloaded: " + err.Error())
	}

	if !pathExists(filepath.Join(profilePath, "chrome")) {
		os.Mkdir(filepath.Join(profilePath, "chrome"), 0644)
	}

	userChromePath := filepath.Join(profilePath, "chrome", "userChrome.css")
	userContentPath := filepath.Join(profilePath, "chrome", "userContent.css")

	if pathExists(userChromePath) {
		os.Rename(userChromePath, userChromePath+".old")
	}

	if pathExists(userContentPath) {
		os.Rename(userContentPath, userContentPath+".old")
	}

	ioutil.WriteFile(userChromePath, userChrome, 0644)
	checkErr(err)
	ioutil.WriteFile(userContentPath, userContent, 0644)
	checkErr(err)

	return nil
}
