package main

import (
	"bufio"
	"os"
	"scripts/utils" // Import the utils package
)

func returnOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getPywalColors(pywalPath string) []string {
	pywalFile, err := os.Open(pywalPath)
	returnOnErr(err)
	defer pywalFile.Close()

	scanner := bufio.NewScanner(pywalFile)
	var allPywalColors []string

	for scanner.Scan() {
		allPywalColors = append(allPywalColors, scanner.Text())
	}
	returnOnErr(scanner.Err())
	return allPywalColors
}

func main() {
	utilArg := os.Args[1]
	wallpaperArg := os.Args[2]
	homeDir, err := os.UserHomeDir()
	returnOnErr(err)
	cavaConfigPath := homeDir + "/.config/cava/config"
	pywalColorsPath := homeDir + "/.cache/wal/colors"
	cachePath := homeDir + "/.cache/wallpaper/"
	var allPywalColors = getPywalColors(pywalColorsPath)

	switch utilArg {
	case "update-wallpaper":
		err = utils.UpdateWallpaper(cachePath, wallpaperArg, cavaConfigPath, allPywalColors)
		returnOnErr(err)
	}
}
