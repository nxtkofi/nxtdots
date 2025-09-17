package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func returnOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	homeDir, err := os.UserHomeDir()
	returnOnErr(err)
	cavaConfigPath := homeDir + "/.config/cava/config"
	pywalColors, error := os.Open(homeDir + "/.cache/wal/colors")
	returnOnErr(error)
	scanner := bufio.NewScanner(pywalColors)

	var allPywalColors []string

	for scanner.Scan() {
		allPywalColors = append(allPywalColors, scanner.Text())
	}

	
}
