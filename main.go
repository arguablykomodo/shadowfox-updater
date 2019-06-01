package main

import (
	"os"
)

var tag string
var header = "Shadowfox updater " + tag

func main() {
	if len(os.Args) > 1 {
		cli()
	} else {
		if err := createUI(); err != nil {
			createFallbackUI()
		}
	}
}
