package main

func main() {
	err := createUI()
	if err != nil {
		createFallbackUI()
	}
}
