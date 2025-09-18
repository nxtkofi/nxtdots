package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func UpdateWallpaper(newWallpaperFullFilePath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cavaConfigPath := homeDir + "/.config/cava/config"
	cachePath := homeDir + "/.cache/wallpaper/"
	useCache := false
	fmt.Println(useCache)

	if _, err := os.Stat(cachePath + "cached_wallpaper"); err == nil {
		useCache = true
	}

	err = os.MkdirAll(filepath.Dir(cachePath+"wallpaper-generated/"), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(cachePath+"current_wallpaper", []byte(newWallpaperFullFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	err = ExecPywal(newWallpaperFullFilePath)
	colors, err := GetPywalColors()

	if err != nil {
		return err
	}
	err = UpdateCavaGradient(cavaConfigPath, colors)
	if err != nil {
		return err
	}
	err = UpdateSpicetify(colors, homeDir)
	// wallpaperFileName := splitWallPaperPath[len(splitWallPaperPath)-1]
	return nil
}
