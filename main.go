package main

const version = "1.5.6"

func main() {
	err := createUI()
	if err != nil {
		createFallbackUI()
	}
}
