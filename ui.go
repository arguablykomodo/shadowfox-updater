package main

import (
	"fmt"

	"github.com/gen2brain/dlgs"
	"github.com/skratchdot/open-golang/open"
)

func createUI() error {
	shouldUpdate, newTag, err := checkForUpdate()
	if err != nil {
		_, err := dlgs.Error(header, err.Error())
		if err != nil {
			return err
		}
	}

	if shouldUpdate {
		wantToUpdate, err := dlgs.Question(header, fmt.Sprintf("There is a new shadowfox-updater version available (%s -> %s)\nDo you want to update?", tag, newTag), true)
		if err != nil {
			return err
		}
		if wantToUpdate {
			open.Start("https://github.com/SrKomodo/shadowfox-updater/#installing")
			return nil
		}
	}

	paths, names, err := getProfilePaths()
	if err != nil {
		_, err := dlgs.Error(header, err.Error())
		if err != nil {
			return err
		}
	}

	name, selected, err := dlgs.List(header, "Which Firefox profile are you going to use?", names)
	if err != nil {
		return err
	}

	if !selected {
		_, err := dlgs.Info(header, "You didn't pick any profile, the application will now close.")
		return err
	}

	pathIndex := 0
	for _, name2 := range names {
		if name == name2 {
			break
		}
		pathIndex++
	}
	profilePath := paths[pathIndex]

	action, selected, err := dlgs.List(header, "What do you want to do?", []string{"Install/Update Shadowfox", "Uninstall Shadowfox"})
	if err != nil {
		return err
	}

	if !selected {
		_, err := dlgs.Info(header, "You didn't pick any action, the application will now close.")
		return err
	}

	if action == "Install/Update Shadowfox" {
		shouldGenerateUUIDs, err := dlgs.Question(header, "Would you like to auto-generate UUIDs?", true)
		if err != nil {
			return err
		}

		shouldSetTheme, err := dlgs.Question(header, "Would you like to automatically set the Firefox dark theme?", false)
		if err != nil {
			return err
		}

		err = install(profilePath, shouldGenerateUUIDs, shouldSetTheme)
		if err == nil {
			_, err = dlgs.Info(header, "Shadowfox has been succesfully installed!")
			if err != nil {
				return err
			}
		} else {
			_, err := dlgs.Error(header, err.Error())
			if err != nil {
				return err
			}
		}
	} else {
		err := uninstall(profilePath)
		if err == nil {
			_, err = dlgs.Info(header, "Shadowfox has been succesfully uninstalled!")
			if err != nil {
				return err
			}
		} else {
			_, err := dlgs.Error(header, err.Error())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
