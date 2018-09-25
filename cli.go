package main

import (
	"flag"
)

func cli() {
	paths, names := getProfilePaths()

	uninstalling := flag.Bool("uninstall", false, "Wheter to install or uninstall ShadowFox")
	profileName := flag.String("profile-name", "", "Name of profile to use, if not defined or not found will fallback to profile-index")
	profileIndex := flag.Int("profile-index", 0, "Index of profile to use")
	uuids := flag.Bool("generate-uuids", false, "Wheter to automatically generate UUIDs or not")
	theme := flag.Bool("set-dark-theme", false, "Wheter to automatically set Firefox's dark theme")

	flag.Parse()

	var path string
	for i, name := range names {
		if name == *profileName {
			path = paths[i]
		}
	}
	if path == "" {
		path = paths[*profileIndex]
	}

	if *uninstalling {
		msg, err := uninstall(path)
		checkErr(msg, err)
	} else {
		msg, err := install(path, *uuids, *theme)
		checkErr(msg, err)
	}
}
