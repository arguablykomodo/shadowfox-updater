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

	themeCheckBox := tview.NewCheckbox().SetLabel("Set Firefox dark theme: ").SetChecked(false)

	installButton := tview.NewButton("Install/Update ShadowFox").SetSelectedFunc(func() {
		message, err := install(paths[profileIndex], uuidCheckBox.IsChecked(), themeCheckBox.IsChecked())
		if err != nil {
			notifyErr(message, err)
		} else {
			infoText.SetText("ShadowFox was successfully installed!")
		}
	})

	uninstallButton := tview.NewButton("Uninstall ShadowFox").SetSelectedFunc(func() {
		message, err := uninstall(paths[profileIndex])
		if err != nil {
			notifyErr(message, err)
		} else {
			infoText.SetText("ShadowFox was successfully uninstalled!")
		}
	})

	exitButton := tview.NewButton("Exit").SetSelectedFunc(func() {
		app.Stop()
	})

	// Setup input callbacks so that TAB switches between buttons
	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyTAB {
			if e.Modifiers()&tcell.ModShift == tcell.ModShift {
				switch {
				case exitButton.HasFocus():
					app.SetFocus(uninstallButton)
				case uninstallButton.HasFocus():
					app.SetFocus(installButton)
				case installButton.HasFocus():
					app.SetFocus(themeCheckBox)
				case themeCheckBox.HasFocus():
					app.SetFocus(uuidCheckBox)
				case uuidCheckBox.HasFocus():
					app.SetFocus(profileSelect)
				case profileSelect.HasFocus():
					app.SetFocus(exitButton)
				}
			} else {
				switch {
				case profileSelect.HasFocus():
					app.SetFocus(uuidCheckBox)
				case uuidCheckBox.HasFocus():
					app.SetFocus(themeCheckBox)
				case themeCheckBox.HasFocus():
					app.SetFocus(installButton)
				case installButton.HasFocus():
					app.SetFocus(uninstallButton)
				case uninstallButton.HasFocus():
					app.SetFocus(exitButton)
				case exitButton.HasFocus():
					app.SetFocus(profileSelect)
				}
			}
		}
		return e
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
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(themeCheckBox, 0, 10, true).
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
	flex.SetBorder(true).
		SetTitle("ShadowFox updater "+version).
		SetBorderPadding(1, 1, 1, 1).
		SetBackgroundColor(tcell.ColorBlack)

	if paths == nil {
		text := tview.NewTextView().SetText(
			"ShadowFox couldn't automatically find 'profiles.ini'. Please follow these steps:\n\n" +
				"1. Close the program                                       \n" +
				"2. Move the program to the folder 'profiles.ini' is located\n" +
				"3. Run the program                                         ",
		).SetTextAlign(tview.AlignCenter)

		text.SetBorder(true).SetTitle("ShadowFox updater "+version).SetBorderPadding(1, 1, 1, 1)

		app.SetRoot(text, true)
	} else {
		app.SetRoot(flex, true).SetFocus(profileSelect)
	}

	if err := app.Run(); err != nil {
		return err
	}
	return nil
}
