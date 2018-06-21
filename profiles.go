package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-ini/ini"
	homedir "github.com/mitchellh/go-homedir"
)

func getProfilePaths() ([]string, []string) {
	// iniPaths stores all profiles.ini files we have to check
	var iniPaths []string

	// Find current directory
	cwd, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Check if profiles.ini exists in the cwd
	exists, _, err := pathExists(filepath.Join(filepath.Dir(cwd), "profiles.ini"))
	if err != nil {
		panic(err)
	}
	if exists { // If it does we just stop here
		iniPaths = []string{filepath.Join(filepath.Dir(cwd), "profiles.ini")}
	} else { // If not we will do some more stuff
		// Get the home directory
		homedir, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		// Possible places where we should check for profiles.ini
		var possible []string

		switch runtime.GOOS {
		case "windows":
			possible = []string{homedir + "\\AppData\\Roaming\\Mozilla\\Firefox\\profiles.ini"}

		case "darwin":
			possible = []string{homedir + "/Library/Application Support/Firefox/profiles.ini"}

		case "linux":
			possible = []string{
				homedir + "/.mozilla/firefox/profiles.ini",
				homedir + "/.mozilla/firefox-trunk/profiles.ini",
			}

		default:
			panic("Sorry, but this program only works on Windows, Mac OS, or Linux")
		}

		// Check if profiles.ini exists on each possible path and add them to the list
		for _, p := range possible {
			exists, _, err := pathExists(p)
			if err != nil {
				panic(err)
			}
			if exists {
				iniPaths = append(iniPaths, p)
				break
			}
		}

		// If we didnt find anything then we just give up
		if len(iniPaths) == 0 {
			return nil, nil
		}
	}

	var paths []string
	var names []string

	// For each possible ini file
	for _, p := range iniPaths {
		file, err := ini.Load(p)
		if err != nil {
			panic(err)
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

	return paths, names
}
