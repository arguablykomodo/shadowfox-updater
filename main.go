package main

var version = "dev"

func main() {
	err := createUI()
	if err != nil {
		createFallbackUI()
	}
}
