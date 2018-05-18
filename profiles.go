package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-ini/ini"
	homedir "github.com/mitchellh/go-homedir"
)

func getProfilePaths() ([]string, []string) {
	var iniPath string

	cwd, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exists, _, err := pathExists(filepath.Join(filepath.Dir(cwd), "profiles.ini"))
	if err != nil {
		panic(err)
	}
	if !exists {
		homedir, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		switch runtime.GOOS {
		case "windows":
			iniPath = homedir + "\\AppData\\Roaming\\Mozilla\\Firefox\\profiles.ini"

		case "darwin":
			iniPath = homedir + "/Library/Application Support/Firefox/profiles.ini"

		case "linux":
			iniPath = homedir + "/.mozilla/firefox/profiles.ini"

		default:
			panic("Sorry, but this program only works on Windows, Mac OS, or Linux")
		}

		exists, _, err := pathExists(iniPath)
		if err != nil {
			panic(err)
		}
		if !exists {
			return nil, nil
		}
	} else {
		iniPath = filepath.Join(filepath.Dir(cwd), "profiles.ini")
	}

	file, err := ini.Load(iniPath)
	if err != nil {
		panic(err)
	}

	var paths []string
	var names []string
	for _, section := range file.Sections() {
		if key, err := section.GetKey("Path"); err == nil {
			path := key.String()
			isRelative := section.Key("IsRelative").MustInt(1)

			if isRelative == 1 {
				paths = append(paths, filepath.Join(filepath.Dir(iniPath), path))
			} else {
				paths = append(paths, path)
			}
			names = append(names, filepath.Base(path))
		}
	}
	return paths, names
}
