package main

import (
	"fmt"
	"strconv"
)

func createFallbackUI() {
	var choice string
	paths, names := getProfilePaths()

	if paths == nil {
		fmt.Println("ShadowFox couldn't automatically find 'profiles.ini'. Please follow these steps:")
		fmt.Println("  1. Close the program")
		fmt.Println("  2. Move the program to the folder 'profiles.ini' is located")
		fmt.Println("  3. Run the program")
		fmt.Scanln()
		return
	}

	fmt.Print("ShadowFox updater 1.4.1\n\n")

	fmt.Println("Available profiles:")
	for i, name := range names {
		fmt.Printf("  %d: %s\n", i, name)
	}

	fmt.Printf("\nWich one would you like to use? [%d-%d] ", 0, len(names)-1)

	var profile string
	for {
		fmt.Scanln(&choice)
		i, err := strconv.Atoi(choice)
		if err != nil || i < 0 || i > len(paths) {
			fmt.Print("Please input a valid number ")
		} else {
			profile = paths[i]
			break
		}
	}

	fmt.Print("\nDo you want to (1) install or (2) uninstall ShadowFox? [1/2] ")
	fmt.Scanln(&choice)

	if choice == "2" {
		uninstall(profile)
		fmt.Print("\nShadowFox was succesfully uninstalled! (Press 'enter' to exit)")
		fmt.Scanln()
		return
	}

	fmt.Print("\nWould you like to auto-generate UUIDs? [y/n] ")
	fmt.Scanln(&choice)

	message, err := install(profile, (choice == "y" || choice == "Y"))
	if err != nil {
		fmt.Printf("%s: %s", message, err.Error())
		fmt.Scanln()
		return
	}

	fmt.Print("\nShadowFox was succesfully installed! (Press 'enter' to exit)")
	fmt.Scanln()
	return
}
