package utils

import (
	"bufio"
	"os"
)

func GetCurrentWallpaperFullPath(homeDir string) []byte {
	fileWithCurrentWallPapaerpath, err := os.Open(homeDir + "/.cache/wallpaper/current_wallpaper")
	ReturnOnErr(err)

	reader := bufio.NewReader(fileWithCurrentWallPapaerpath)
	pathToCurrentWallpaper, _, err := reader.ReadLine()
	ReturnOnErr(err)

	return pathToCurrentWallpaper
}
