package main

import (
	"os"

	"github.com/gen2brain/dlgs"
)

func checkErr(msg string, err error) {
	if err != nil {
		dlgs.Error("Shadowfox Updater", msg+"\n"+err.Error())
		panic(err)
	}
}

func main() {
	if len(os.Args) > 1 {
		cli()
	} else {
		err := createUI()
		if err != nil {
			createFallbackUI()
		}
	}
}
