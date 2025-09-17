package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func UpdateWallpaper(wallpaperCachePath string, newWallpaperFullFilePath string, cavaConfigPath string, allPywalColors []string) error {
	useCache := false
	fmt.Println(useCache)
	if _, err := os.Stat(wallpaperCachePath + "cached_wallpaper"); err == nil {
		useCache = true
	}

	err := os.MkdirAll(filepath.Dir(wallpaperCachePath+"wallpaper-generated/"), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(wallpaperCachePath+"current_wallpaper", []byte(newWallpaperFullFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	err = Pywal(newWallpaperFullFilePath)
	if err != nil {
		return err
	}
	err = UpdateCavaGradient(cavaConfigPath, allPywalColors)
	if err != nil {
		return err
	}
	err = UpdateSpicetify()
	// wallpaperFileName := splitWallPaperPath[len(splitWallPaperPath)-1]
	return nil
}
