package main

const version = "1.5.5"

func main() {
	err := createUI()
	if err != nil {
		createFallbackUI()
	}
}
