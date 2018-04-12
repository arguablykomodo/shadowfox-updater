package main

import (
	"github.com/marcusolsson/tui-go"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func center(widget tui.Widget) *tui.Box {
	return tui.NewHBox(
		tui.NewSpacer(),
		widget,
		tui.NewSpacer(),
	)
}
