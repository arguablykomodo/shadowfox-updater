package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	homedir "github.com/mitchellh/go-homedir"
)

func getProfilePaths() ([]string, []string, error) {
	// iniPaths stores all profiles.ini files we have to check
	iniPaths := []string{}

	// Get the home directory
	homedir, err := homedir.Dir()
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't find home directory: %s", err)
	}

	// Possible places where we should check for profiles.ini
	possible := []string{
		"./profiles.ini",
		homedir + "\\AppData\\Roaming\\Mozilla\\Firefox\\profiles.ini",
		homedir + "/Library/Application Support/Firefox/profiles.ini",
		homedir + "/.mozilla/firefox/profiles.ini",
		homedir + "/.mozilla/firefox-trunk/profiles.ini",
	}

	// Check if profiles.ini exists on each possible path and add them to the list
	for _, p := range possible {
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, nil, fmt.Errorf("Couldn't check if %s exists: %s", p, err)
		}
		iniPaths = append(iniPaths, p)
	}

	// If we didnt find anything then we just give up
	if len(iniPaths) == 0 {
		return nil, nil, errors.New("Couldn't find any profiles")
	}

	var paths []string
	var names []string

	// For each possible ini file
	for _, p := range iniPaths {
		file, err := ini.Load(p)
		if err != nil {
			return nil, nil, fmt.Errorf("Could not read profiles.ini, make sure its encoded in UTF-8: %s", err)
		}

		// Find the Path key and add it to the list
		for _, section := range file.Sections() {
			if key, err := section.GetKey("Path"); err == nil {
				path := key.String()
				isRelative := section.Key("IsRelative").MustInt(1)

				if isRelative == 1 {
					paths = append(paths, filepath.Join(filepath.Dir(p), path))
				} else {
					paths = append(paths, path)
				}
				names = append(names, filepath.Base(path))
			}
		}
	}

	return paths, names, nil
}
