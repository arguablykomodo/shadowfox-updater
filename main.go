package main

import (
	"os"
)

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
