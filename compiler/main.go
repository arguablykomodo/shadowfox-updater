package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func writeManifest() {
	findStr := []byte(`name="SrKomodo.Software.shadowfoxUpdater"`)

	err := os.Remove("manifest.xml")
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	manifest, err := ioutil.ReadFile("_manifest.xml")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(
		"manifest.xml",
		bytes.Replace(
			manifest,
			findStr,
			append(
				findStr,
				[]byte(" version=\""+strings.TrimPrefix(os.Args[1], "v")+".0\"")...,
			),
			1,
		),
		0644,
	)
	if err != nil {
		panic(err)
	}
}

func main() {
	version := strings.Split(strings.TrimPrefix(os.Args[1], "v"), ".")
	major := version[0]
	minor := version[1]
	patch := version[2]

	var err error

	osList := [3]string{"windows", "darwin", "linux"}
	archList := [2]string{"386", "amd64"}

	osNames := [3]string{"windows", "mac", "linux"}
	archNames := [2]string{"x32", "x64"}

	writeManifest()

	// Generate .syso files
	fmt.Printf("Generating .syso files\n")
	sysoArgs := []string{
		"-platform-specific=true",
		"-ver-major=" + major,
		"-ver-minor=" + minor,
		"-ver-patch=" + patch,
		"-ver-build=0",
		"-product-ver-major=" + major,
		"-product-ver-minor=" + minor,
		"-product-ver-patch=" + patch,
		"-product-ver-build=0",
	}
	err = exec.Command("goversioninfo", sysoArgs...).Run()
	checkErr(err)

	for i, buildOS := range osList {
		for j, buildArch := range archList {
			// Set env variables
			err = os.Setenv("GOOS", buildOS)
			checkErr(err)
			err = os.Setenv("GOARCH", buildArch)
			checkErr(err)

			args := []string{
				"build",
				"-ldflags", `"-X main.version=` + os.Args[1] + `"`,
				"-o", "dist/shadowfox_" + osNames[i] + "_" + archNames[j],
			}

			if buildOS == "windows" {
				args[4] += ".exe"
			}

			fmt.Printf("go %v\n", strings.Join(args, " "))

			fmt.Printf("Compiling %v %v\n", buildOS, buildArch)
			err = exec.Command("go", args...).Run()
			checkErr(err)
		}
	}
}
