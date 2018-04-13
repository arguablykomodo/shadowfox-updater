package main

import (
	"github.com/marcusolsson/tui-go"
)

func center(widget tui.Widget) *tui.Box {
	return tui.NewHBox(
		tui.NewSpacer(),
		widget,
		tui.NewSpacer(),
	)
}

var infoLabel *tui.Label

func createUI() *tui.UI {
	theme := tui.DefaultTheme
	theme.SetStyle("label.error", tui.Style{Bg: tui.ColorRed, Fg: tui.ColorBlack})

	infoLabel := tui.NewLabel("")

	installButton := tui.NewButton("[Install or update ShadowFox]")
	uninstallButton := tui.NewButton("[Uninstall ShadowFox]")

	container := tui.NewVBox(tui.NewPadder(1, 0, tui.NewVBox(
		center(tui.NewLabel("Welcome to the ShadowFox updater 1.0.0")),
		center(tui.NewLabel("What would you like to do?")),
		tui.NewPadder(0, 1, tui.NewVBox(
			center(installButton),
			center(uninstallButton),
		)))),
		center(infoLabel),
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

	ui.SetTheme(theme)
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	tui.DefaultFocusChain.Set(installButton, uninstallButton)

	installButton.OnActivated(func(b *tui.Button) {
		err := install()
		if err != nil {
			infoLabel.SetStyleName("error")
			infoLabel.SetText(err.Error())
		}
		infoLabel.SetText("ShadowFox has been succesfully installed/updated")
	})

	return &ui
}
