package main

import (
	"os"
)

var version = "dev"

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
