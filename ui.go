package main

import (
	"github.com/marcusolsson/tui-go"
)

func createUI() *tui.UI {
	installButton := tui.NewButton("[Install ShadowFox]")
	updateButton := tui.NewButton("[Upgrade ShadowFox]")
	uninstallButton := tui.NewButton("[Uninstall]")

	container := tui.NewVBox(tui.NewPadder(1, 0, tui.NewVBox(
		tui.NewLabel("Welcome to the ShadowFox updater 1.0.0"),
		tui.NewPadder(6, 0, tui.NewLabel("What would you like to do?")),
		tui.NewPadder(1, 1, tui.NewVBox(
			center(installButton),
			center(updateButton),
			center(uninstallButton),
		)))),
	)
	container.SetBorder(true)

	root := tui.NewVBox(
		tui.NewSpacer(),
		tui.NewHBox(
			tui.NewSpacer(),
			container,
			tui.NewSpacer(),
		),
		tui.NewSpacer(),
	)

	ui, err := tui.New(root)
	checkErr(err)

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	tui.DefaultFocusChain.Set(installButton, updateButton, uninstallButton)

	return &ui
}
