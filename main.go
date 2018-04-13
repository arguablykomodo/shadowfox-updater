package main

import (
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"

	"github.com/go-ini/ini"
)

var profilePath string

func main() {
	profilePath = getProfilePath()

	err := (*createUI()).Run()
	checkErr(err)
}

func getProfilePath() string {
	var iniPath string

	homedir, err := homedir.Dir()
	checkErr(err)

	switch runtime.GOOS {
	case "windows":
		iniPath = homedir + "\\AppData\\Roaming\\Mozilla\\Firefox\\profiles.ini"

	case "darwin":
		iniPath = homedir + "/Library/Mozilla/profiles.ini"

	case "linux":
		iniPath = homedir + "/.mozilla/profiles.ini"

	default:
		panic("Sorry, but this program only works on Windows, Mac OS, or Linux")
	}

	file, err := ini.Load(iniPath)
	checkErr(err)

	for _, section := range file.Sections() {
		if section.Key("Default").MustInt(0) == 1 {
			path := section.Key("Path").MustString("")
			isRelative := section.Key("IsRelative").MustInt(1)

			if isRelative == 1 {
				return filepath.Join(filepath.Dir(iniPath), path)
			}
			return path
		}
	}

	panic("The default profile couln't be found in profiles.ini")
}
