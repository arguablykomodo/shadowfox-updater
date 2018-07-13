package main

const version = "1.5.3"

func main() {
	err := createUI()
	if err != nil {
		createFallbackUI()
	}
}
