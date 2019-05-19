package main

import (
	"github.com/gen2brain/dlgs"
)

func createUI() error {
	paths, names := getProfilePaths()

	name, selected, err := dlgs.List("Shadowfox Updater", "Which Firefox profile are you going to use?", names)
	checkErr("", err)
	if !selected {
		dlgs.Info("Shadowfox Updater", "You didn't pick any profile, the application will now close.")
		return nil
	}

	pathIndex := 0
	for _, name2 := range names {
		if name == name2 {
			break
		}
		pathIndex++
	}
	profilePath := paths[pathIndex]

	action, selected, err := dlgs.List("Shadowfox Updater", "What do you want to do?", []string{"Install/Update Shadowfox", "Uninstall Shadowfox"})
	checkErr("", err)
	if !selected {
		dlgs.Info("Shadowfox Updater", "You didn't pick any action, the application will now close.")
		return nil
	}

	if action == "Install/Update Shadowfox" {
		shouldGenerateUUIDs, err := dlgs.Question("Shadowfox Updater", "Would you like to auto-generate UUIDs?", true)
		checkErr("", err)

		shouldSetTheme, err := dlgs.Question("Shadowfox Updater", "Would you like to automatically set the Firefox dark theme?", false)
		checkErr("", err)

		msg, err := install(profilePath, shouldGenerateUUIDs, shouldSetTheme)
		if err == nil {
			dlgs.Info("Shadowfox Updater", "Shadowfox has been succesfully installed!")
		} else {
			checkErr(msg, err)
		}
	} else {
		msg, err := uninstall(profilePath)
		if err == nil {
			dlgs.Info("Shadowfox Updater", "Shadowfox has been succesfully uninstalled!")
		} else {
			checkErr(msg, err)
		}
	}
	return nil
}
