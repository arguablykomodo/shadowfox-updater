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

func createUI() {
	app := tview.NewApplication()
	paths := getProfilePaths()
	profileIndex := 0

	infoText = tview.NewTextView()
	infoText.SetText("(Press TAB for selecting, ENTER for clicking)")
	infoText.SetTextAlign(tview.AlignCenter)

	// Create buttons
	profileSelect := tview.NewDropDown().SetLabel("Profile to use: ").SetOptions(paths, func(text string, index int) {
		profileIndex = index
	})

	installButton := tview.NewButton("Install/Update ShadowFox").SetSelectedFunc(func() {
		message, err := install(paths[profileIndex])
		if err != nil {
			notifyErr(message, err)
		} else {
			infoText.SetText("ShadowFox was succesfully installed!")
		}
	})

	uninstallButton := tview.NewButton("Uninstall ShadowFox").SetSelectedFunc(func() {
		err := os.RemoveAll(filepath.Join(paths[profileIndex], "chrome"))
		if err != nil {
			notifyErr("Couln't uninstall ShadowFox", err)
		} else {
			infoText.SetText("Shadowfox was succesfully uninstalled!")
		}
	})

	exitButton := tview.NewButton("Exit").SetSelectedFunc(func() {
		app.Stop()
	})

	// Setup input so that TAB switches between buttons
	profileSelect.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			app.SetFocus(installButton)
		}
		return event
	})

	installButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			app.SetFocus(uninstallButton)
		}
		return event
	})

	uninstallButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			app.SetFocus(exitButton)
		}
		return event
	})

	exitButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			app.SetFocus(profileSelect)
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
		AddItem(infoText, 5, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(installButton, 0, 1, true).
			AddItem(nil, 1, 0, false).
			AddItem(uninstallButton, 0, 1, true).
			AddItem(nil, 1, 0, false).
			AddItem(exitButton, 0, 1, true), 3, 0, true,
		)

	flex.SetBorder(true).SetTitle("ShadowFox updater 1.0.0").SetBorderPadding(1, 1, 1, 1)
	if err := app.SetRoot(flex, true).SetFocus(profileSelect).Run(); err != nil {
		panic(err)
	}
}

func getProfilePaths() []string {
	var iniPath string

	homedir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

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
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, section := range file.Sections() {
		if key, err := section.GetKey("Path"); err == nil {
			path := key.String()
			isRelative := section.Key("IsRelative").MustInt(1)

			if isRelative == 1 {
				paths = append(paths, filepath.Join(filepath.Dir(iniPath), path))
			} else {
				paths = append(paths, path)
			}
		}
	}
	return paths
}
