package main

import (
	"github.com/gdamore/tcell"

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

func createUI() error {
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
	flex.SetBorder(true).SetTitle("ShadowFox updater 1.5.1").SetBorderPadding(1, 1, 1, 1)

	if paths == nil {
		text := tview.NewTextView().SetText(
			"ShadowFox couldn't automatically find 'profiles.ini'. Please follow these steps:\n\n" +
				"1. Close the program                                       \n" +
				"2. Move the program to the folder 'profiles.ini' is located\n" +
				"3. Run the program                                         ",
		).SetTextAlign(tview.AlignCenter)

		text.SetBorder(true).SetTitle("ShadowFox updater 1.5.1").SetBorderPadding(1, 1, 1, 1)

		app.SetRoot(text, true)
	} else {
		app.SetRoot(flex, true).SetFocus(profileSelect)
	}

	if err := app.Run(); err != nil {
		return err
	}
	return nil
}
