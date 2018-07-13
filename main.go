package main

const version = "1.5.4"

func main() {
	err := createUI()
	if err != nil {
		createFallbackUI()
	}
}
