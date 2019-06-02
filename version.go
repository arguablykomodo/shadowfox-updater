package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

var tag string

func checkForUpdate() (bool, string, error) {
	resp, err := http.Get("https://api.github.com/repos/SrKomodo/shadowfox-updater/releases/latest")
	if err != nil {
		return false, "", err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	regex := regexp.MustCompile(`"tag_name":"(.+?)"`)
	newTag := string(regex.FindSubmatch(data)[1])

	return newTag != tag, newTag, nil
}
