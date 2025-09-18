package main

import (
	"os"
	"scripts/utils"
)

func returnOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	utilArg := os.Args[1]
	wallpaperArg := os.Args[2]

	switch utilArg {
	case "update-wallpaper":
		err := utils.UpdateWallpaper(wallpaperArg)
		returnOnErr(err)
	case "update-spicetify":
		homeDir, err := os.UserHomeDir()
		returnOnErr(err)

		colors, err := utils.GetPywalColors()
		returnOnErr(err)

		err = utils.UpdateSpicetify(colors, homeDir)

	}
}
