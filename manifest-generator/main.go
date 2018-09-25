package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
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
			-1,
		),
		0644,
	)
	if err != nil {
		panic(err)
	}
}
