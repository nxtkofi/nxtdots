package main

import (
	"os"
	"scripts/utils"
)

func main() {
	utilArg := os.Args[1]
	homeDir, err := os.UserHomeDir()
	utils.ReturnOnErr(err)

	switch utilArg {
	case "update-wallpaper":
		utils.UpdateWallpaper(os.Args[2], homeDir)
		utils.ReturnOnErr(err)
	case "system-theme":
		err := utils.HandleThemeChange(os.Args[2], homeDir)
		utils.ReturnOnErr(err)
	case "install":
		utils.Install()
	case "rice-spotify":
		utils.RiceSpotify()
	case "restart-waybar":
		err := utils.RestartWaybar()
		utils.ReturnOnErr(err)
	}
}
