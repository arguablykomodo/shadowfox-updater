package main

import (
	"os"
	"os/exec"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error

	osList := []string{"windows", "darwin", "linux"}
	archList := []string{"386", "amd64"}

	for _, buildOS := range osList {
		for _, buildArch := range archList {
			// Set env variables
			err = os.Setenv("GOOS", buildOS)
			checkErr(err)
			err = os.Setenv("GOARCH", buildArch)
			checkErr(err)

			args := []string{"build", "-o", "dist/shadowfox_" + buildOS + "_" + buildArch}

			if buildOS == "windows" {
				args[2] += ".exe"

				// Generate .syso files
				rsrc := exec.Command(
					"rsrc",
					"-manifest", "manifest.xml",
					"-arch", buildArch,
					"-o", "shadowfox.syso",
				)
				_, err := rsrc.Output()
				checkErr(err)
			}

			build := exec.Command("go", args...)
			_, err := build.Output()
			checkErr(err)
		}
	}
}
