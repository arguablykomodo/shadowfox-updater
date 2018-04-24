package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/gdamore/tcell"

	"github.com/go-ini/ini"
	"github.com/mitchellh/go-homedir"
	"github.com/rivo/tview"
)

var infoText *tview.TextView

func notifyErr(message string, err error) {
	infoText.SetTextColor(tcell.ColorWhite)
	infoText.SetBackgroundColor(tcell.ColorRed)
	if message != "" {
		infoText.SetText(message + ":\n" + err.Error())
	} else {
		infoText.SetText(err.Error())
	}
}

func uninstall(profile string) (string, error) {
	err := os.RemoveAll(filepath.Join(profile, "chrome", "ShadowFox_customization"))
	if err != nil {
		return "Couldn't delete ShadowFox_customization", err
	}
	err = os.Remove(filepath.Join(profile, "chrome", "userChrome.css"))
	if err != nil {
		return "Couln't delete userChrome.css", err
	}
	err = os.Remove(filepath.Join(profile, "chrome", "userContent.css"))
	if err != nil {
		return "Couln't delete userContent.css", err
	}
	return "", nil
}

func createUI() {
	app := tview.NewApplication()
	paths, names := getProfilePaths()
	profileIndex := 0

	infoText = tview.NewTextView()
	infoText.SetText("(Press TAB for selecting, ENTER for clicking)")
	infoText.SetTextAlign(tview.AlignCenter)

	// Create buttons
	profileSelect := tview.NewDropDown().SetLabel("Profile to use: ").SetOptions(names, func(text string, index int) {
		profileIndex = index
	})

	uuidCheckBox := tview.NewCheckbox().SetLabel("Auto-Generate UUIDs: ").SetChecked(false)

	installButton := tview.NewButton("Install/Update ShadowFox").SetSelectedFunc(func() {
		message, err := install(paths[profileIndex], uuidCheckBox.IsChecked())
		if err != nil {
			notifyErr(message, err)
		} else {
			infoText.SetText("ShadowFox was succesfully installed!")
		}
	})

	uninstallButton := tview.NewButton("Uninstall ShadowFox").SetSelectedFunc(func() {
		message, err := uninstall(paths[profileIndex])
		if err != nil {
			notifyErr(message, err)
		} else {
			infoText.SetText("ShadowFox was succesfully uninstalled!")
		}
	})

	exitButton := tview.NewButton("Exit").SetSelectedFunc(func() {
		app.Stop()
	})

	// Setup input so that TAB switches between buttons
	profileSelect.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			if (event.Modifiers() & tcell.ModShift) == tcell.ModShift {
				app.SetFocus(exitButton)
			} else {
				app.SetFocus(uuidCheckBox)
			}
		}
		return event
	})

	uuidCheckBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			if (event.Modifiers() & tcell.ModShift) == tcell.ModShift {
				app.SetFocus(profileSelect)
			} else {
				app.SetFocus(installButton)
			}
		}
		return event
	})

	installButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			if (event.Modifiers() & tcell.ModShift) == tcell.ModShift {
				app.SetFocus(uuidCheckBox)
			} else {
				app.SetFocus(uninstallButton)
			}
		}
		return event
	})

	uninstallButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			if (event.Modifiers() & tcell.ModShift) == tcell.ModShift {
				app.SetFocus(installButton)
			} else {
				app.SetFocus(exitButton)
			}
		}
		return event
	})

	exitButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			if (event.Modifiers() & tcell.ModShift) == tcell.ModShift {
				app.SetFocus(uninstallButton)
			} else {
				app.SetFocus(profileSelect)
			}
		}
		return event
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(profileSelect, 0, 10, true).
			AddItem(nil, 0, 1, false), 1, 0, true,
		).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(uuidCheckBox, 0, 10, true).
			AddItem(nil, 0, 1, false), 1, 0, true,
		).
		AddItem(nil, 1, 0, false).
		AddItem(infoText, 5, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(installButton, 0, 1, true).
			AddItem(nil, 1, 0, false).
			AddItem(uninstallButton, 0, 1, true).
			AddItem(nil, 1, 0, false).
			AddItem(exitButton, 0, 1, true), 3, 0, true,
		)
	flex.SetBorder(true).SetTitle("ShadowFox updater 1.4.1").SetBorderPadding(1, 1, 1, 1)

	if paths == nil {
		text := tview.NewTextView().SetText(
			"ShadowFox couldn't automatically find 'profiles.ini'. Please follow these steps:\n\n" +
				"1. Close the program                                       \n" +
				"2. Move the program to the folder 'profiles.ini' is located\n" +
				"3. Run the program                                         ",
		).SetTextAlign(tview.AlignCenter)

		text.SetBorder(true).SetTitle("ShadowFox updater 1.4.1").SetBorderPadding(1, 1, 1, 1)

		app.SetRoot(text, true)
	} else {
		app.SetRoot(flex, true).SetFocus(profileSelect)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}

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
