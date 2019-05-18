package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

func createUI() error {
	err := ui.Main(func() {
		paths, names := getProfilePaths()

		infoText := ui.NewLabel("")

		dropdown := ui.NewCombobox()
		for _, name := range names {
			dropdown.Append(name)
		}
		dropdown.SetSelected(0)

		dropdownForm := ui.NewForm()
		dropdownForm.SetPadded(true)
		dropdownForm.Append("Profile to use:", dropdown, false)

		installButton := ui.NewButton("Install/Update ShadowFox")
		uninstallButton := ui.NewButton("Uninstall ShadowFox")
		exitButton := ui.NewButton("Exit")

		buttons := ui.NewHorizontalBox()
		buttons.SetPadded(true)
		buttons.Append(installButton, true)
		buttons.Append(uninstallButton, true)
		buttons.Append(exitButton, true)

		uuidCheckBox := ui.NewCheckbox("Auto-Generate UUIDs")
		themeCheckBox := ui.NewCheckbox("Set Firefox dark theme")

		installButton.OnClicked(func(b *ui.Button) {
			message, err := install(paths[dropdown.Selected()], uuidCheckBox.Checked(), themeCheckBox.Checked())
			if err != nil {
				infoText.SetText(message + ":\n" + err.Error())
			} else {
				infoText.SetText("ShadowFox was successfully installed!")
			}
		})

		uninstallButton.OnClicked(func(b *ui.Button) {
			message, err := uninstall(paths[dropdown.Selected()])
			if err != nil {
				infoText.SetText(message + ":\n" + err.Error())
			} else {
				infoText.SetText("ShadowFox was successfully uninstalled!")
			}
		})

		exitButton.OnClicked(func(b *ui.Button) {
			ui.Quit()
		})

		box := ui.NewVerticalBox()
		box.SetPadded(true)
		box.Append(infoText, false)
		box.Append(dropdownForm, false)
		box.Append(uuidCheckBox, false)
		box.Append(themeCheckBox, false)
		box.Append(buttons, false)

		window := ui.NewWindow("Shadowfox Updater", 1, 1, false)
		window.SetMargined(true)
		window.SetChild(box)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	return err
}
