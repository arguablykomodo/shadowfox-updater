package main

import (
	"flag"
	"fmt"
)

func cli() {
	paths, names, err := getProfilePaths()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Scanln()
		panic(err)
	}

	version := flag.Bool("version", false, "Shows the current version")
	uninstalling := flag.Bool("uninstall", false, "Wheter to install or uninstall ShadowFox")
	profileName := flag.String("profile-name", "", "Name of profile to use, if not defined or not found will fallback to profile-index")
	profileIndex := flag.Int("profile-index", 0, "Index of profile to use")
	uuids := flag.Bool("generate-uuids", false, "Wheter to automatically generate UUIDs or not")
	theme := flag.Bool("set-dark-theme", false, "Wheter to automatically set Firefox's dark theme")

	flag.Parse()

	if *version {
		fmt.Println(header)
		return
	}

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
		err := uninstall(path)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Scanln()
			panic(err)
		}
	} else {
		err := install(path, *uuids, *theme)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Scanln()
			panic(err)
		}
	}
}
