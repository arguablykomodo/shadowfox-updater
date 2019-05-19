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
	osList := [3]string{"windows", "darwin", "linux"}
	archList := [2]string{"386", "amd64"}

	osNames := [3]string{"windows", "mac", "linux"}
	archNames := [2]string{"x32", "x64"}

	var err error
	for i, buildOS := range osList {
		for j, buildArch := range archList {
			err = os.Setenv("GOOS", buildOS)
			checkErr(err)
			err = os.Setenv("GOARCH", buildArch)
			checkErr(err)

			args := []string{"build", "-o", "dist/shadowfox_" + osNames[i] + "_" + archNames[j]}

			if buildOS == "windows" {
				args[2] += ".exe"
			}

			_, err := exec.Command("go", args...).Output()
			checkErr(err)
		}
	}
}
